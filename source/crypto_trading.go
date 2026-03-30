package main

import (
	"fmt"
	"math"
	"strings"
)

// TradingConfig holds the thresholds for buy/sell signals
type TradingConfig struct {
	Oversold   int     // RSI threshold for oversold (default: 30)
	Overbought int     // RSI threshold for overbought (default: 70)
	BBStdDev   float64 // Bollinger Bands standard deviation (default: 2.0)
}

// TradingSignal represents the action to take for a coin
type TradingSignal struct {
	Action string // BUY, SELL, HOLD
	Reason string // Why this action (RSI oversold, below BB, etc.)
	Target string // Which coin to rotate to (if SELL)
	Units  float64
	Price  float64
}

// DefaultTradingConfig returns the default configuration
func DefaultTradingConfig() TradingConfig {
	return TradingConfig{
		Oversold:   30,
		Overbought: 70,
		BBStdDev:   2.0,
	}
}

// CalculateTradingSignal determines the trading action for a coin
// Logic:
// 1. If RSI <= oversold OR price below lower BB → BUY (this coin holds, others rotate to it)
// 2. If RSI >= overbought OR price above upper BB → SELL (sell from this coin into lower RSI coin)
// 3. If nothing in oversold/overbought range → HOLD
// 4. If current coin overbought but NO oversold coins exist → sell into USD
func CalculateTradingSignal(sym string, currentPrice, entryPrice, held float64, item CryptoItem, items []CryptoItem, config TradingConfig) string {
	// Skip USD/USDC
	if sym == "USD" || sym == "USDC" {
		return "HOLD"
	}

	// Determine if current coin is oversold or overbought based on RSI
	isOversold := item.RSI > 0 && item.RSI <= float64(config.Oversold)
	isOverbought := item.RSI > 0 && item.RSI >= float64(config.Overbought)

	// Check Bollinger Bands
	isBelowBBL := item.BBLower > 0 && currentPrice < item.BBLower
	isAboveBBU := item.BBUpper > 0 && currentPrice > item.BBUpper

	// Combined signals - use both RSI and BB
	isBuySignal := isOversold || isBelowBBL
	isSellSignal := isOverbought || isAboveBBU

	// Find the absolute lowest RSI coin among ALL coins (not just other coins)
	absoluteLowestRSICoin := ""
	absoluteLowestRSI := 101.0
	for _, it := range items {
		if it.Sym == "USD" || it.Sym == "USDC" {
			continue
		}
		if it.RSI > 0 && it.RSI < absoluteLowestRSI {
			absoluteLowestRSI = it.RSI
			absoluteLowestRSICoin = it.Sym
		}
	}

	// Case 1: Current coin is oversold (RSI < oversold OR below BB lower) → HOLD
	// Others should rotate INTO this coin
	if isBuySignal {
		reason := "RSI oversold"
		if isBelowBBL {
			reason = "Below BB lower"
		}
		return fmt.Sprintf("HOLD (%s)", reason)
	}

	// Case 2: Current coin is NOT overbought → HOLD
	// No sell signal means we don't sell
	if !isSellSignal {
		return "HOLD"
	}

	// Case 3: Current coin is overbought (RSI > overbought OR above BB upper)
	// If current coin has lowest RSI → HOLD (best buy, don't sell)
	isLowestRSI := item.RSI > 0 && item.RSI <= absoluteLowestRSI

	if isLowestRSI {
		return "HOLD"
	}

	// Determine rotation target
	rotation := "USD"

	// If there's a coin with lower RSI (that isn't us), rotate to it
	if absoluteLowestRSICoin != "" && absoluteLowestRSICoin != sym {
		rotation = absoluteLowestRSICoin
	}
	// Otherwise rotate to USD (no better coin to rotate into)

	// Calculate units to sell (25% of position)
	var unitsToSell float64
	if held > 0 {
		if held < 1.0 {
			unitsToSell = held * 0.25
			if unitsToSell < 0.01 {
				unitsToSell = held
			}
		} else {
			unitsToSell = math.Floor(held * 0.25)
			if unitsToSell < 1 {
				unitsToSell = 1
			}
		}
	}

	// Calculate target price based on profit
	profitPct := ((currentPrice - entryPrice) / entryPrice) * 100
	var target float64
	if profitPct >= 50 {
		target = currentPrice * 1.15
	} else if profitPct >= 20 {
		target = currentPrice * 1.20
	} else {
		target = currentPrice * 1.25
	}

	// Format units
	var unitsStr string
	if unitsToSell < 0.1 {
		unitsStr = fmt.Sprintf("%.2f", unitsToSell)
	} else if unitsToSell < 1 {
		unitsStr = fmt.Sprintf("%.1f", unitsToSell)
	} else {
		unitsStr = fmt.Sprintf("%.0f", unitsToSell)
	}

	// Format target price with comma separators if > 999
	targetStr := fmt.Sprintf("%d", int(target))
	if len(targetStr) > 3 {
		// Add comma separators
		var result strings.Builder
		length := len(targetStr)
		for i, c := range targetStr {
			if i > 0 && (length-i)%3 == 0 {
				result.WriteString(",")
			}
			result.WriteRune(c)
		}
		targetStr = result.String()
	}

	return fmt.Sprintf("%s @ $%s → %s", unitsStr, targetStr, rotation)
}

// GetBestBuyTarget returns the coin with the best buy opportunity (lowest RSI)
// This is used for USD to suggest which coin to buy
func GetBestBuyTarget(items []CryptoItem, config TradingConfig) string {
	bestBuy := ""
	bestRSI := 101.0 // RSI is 0-100, use 101 as sentinel

	for _, it := range items {
		if it.Sym == "USD" || it.Sym == "USDC" {
			continue
		}
		// Look for oversold coins
		isOversold := it.RSI > 0 && it.RSI <= float64(config.Oversold)
		if isOversold && it.RSI < bestRSI {
			bestRSI = it.RSI
			bestBuy = it.Sym
		}
	}

	return bestBuy
}
