package crypto

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

// formatUnits formats the held/units amount
func formatUnits(units float64) string {
	if units < 0.1 {
		return fmt.Sprintf("%.2f", units)
	} else if units < 1 {
		return fmt.Sprintf("%.1f", units)
	}
	return fmt.Sprintf("%.0f", units)
}

// formatPrice formats a price with comma separators if > 999
func formatPrice(price float64) string {
	priceStr := fmt.Sprintf("%.0f", price)
	if len(priceStr) > 3 {
		var result strings.Builder
		length := len(priceStr)
		for i, c := range priceStr {
			if i > 0 && (length-i)%3 == 0 {
				result.WriteString(",")
			}
			result.WriteRune(c)
		}
		priceStr = result.String()
	}
	return priceStr
}

// CalculateTradingSignal determines the trading action for a coin
// Logic:
// 1. If RSI <= oversold OR price below lower BB → BUY (this coin holds, others rotate to it)
// 2. If RSI >= overbought OR price above upper BB → SELL (sell from this coin into lower RSI coin)
// 3. If nothing in oversold/overbought range → HOLD
// 4. If current coin overbought but NO oversold coins exist → sell into USD
func CalculateTradingSignal(sym string, currentPrice, entryPrice, held float64, item CryptoItem, items []CryptoItem, config TradingConfig) string {
	// Handle USD/USDC - show best buy target
	if sym == "USD" || sym == "USDC" {
		if held > 0 {
			// Has funds to rebalance - show best buy target
			bestBuy := GetBestBuyTarget(items, config)
			if bestBuy != "" {
				return fmt.Sprintf("→ %s", bestBuy)
			}
		}
		return ""
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

	// Calculate units and target price (reuse for both SELL and HOLD)
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

	// Calculate target price based on 90-day high (10% below high)
	target := currentPrice * 1.25 // default fallback
	if len(item.OHLCData) > 0 {
		high := item.OHLCData[0]
		for _, p := range item.OHLCData {
			if p > high {
				high = p
			}
		}
		target = high * 0.85
	}

	// Format the sell string (units @ price)
	unitsStr := formatUnits(unitsToSell)
	targetStr := formatPrice(target)
	sellStr := fmt.Sprintf("%s @ $%s", unitsStr, targetStr)

	// Case 1: Current coin is oversold (RSI < oversold OR below BB lower) → HOLD
	// Others should rotate INTO this coin
	if isBuySignal {
		reason := "RSI oversold"
		if isBelowBBL {
			reason = "Lower BB"
		}
		return fmt.Sprintf("%s (%s)", sellStr, reason)
	}

	// Case 2: Current coin is NOT overbought → HOLD
	// No sell signal means we don't sell
	if !isSellSignal {
		return fmt.Sprintf("%s (Limit)", sellStr)
	}

	// Case 3: Current coin is overbought (RSI > overbought OR above BB upper)
	// If current coin has lowest RSI → HOLD (best buy, don't sell)
	isLowestRSI := item.RSI > 0 && item.RSI <= absoluteLowestRSI

	if isLowestRSI {
		return fmt.Sprintf("%s (Limit)", sellStr)
	}

	// Determine rotation target and reason
	rotation := "USD"
	reason := "RSI overbought"
	if isAboveBBU {
		reason = "Upper BB"
	}

	// If there's a coin with lower RSI (that isn't us), rotate to it
	if absoluteLowestRSICoin != "" && absoluteLowestRSICoin != sym {
		rotation = absoluteLowestRSICoin
	}
	// Otherwise rotate to USD (no better coin to rotate into)

	return fmt.Sprintf("%s → %s (%s)", sellStr, rotation, reason)
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
