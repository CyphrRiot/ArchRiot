package crypto

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Config represents the crypto.toml structure
type Config struct {
	Indicators IndicatorsConfig `toml:"indicators"`
	Pairs      []PairConfig     `toml:"pairs"`
	Display    DisplaySettings  `toml:"display"`
	APIKey     string           `toml:"api_key"`
}

type PairConfig struct {
	Sym   string  `toml:"sym"`
	Coin  string  `toml:"coin"`
	Held  float64 `toml:"held"`
	Entry float64 `toml:"entry"`
}

type IndicatorsConfig struct {
	RSIPeriod  int     `toml:"rsi_period"`
	Oversold   int     `toml:"oversold"`
	Overbought int     `toml:"overbought"`
	BBPeriod   int     `toml:"bb_period"`
	BBStdDev   float64 `toml:"bb_std"`
}

type DisplaySettings struct {
	ShowTotals bool `toml:"show_totals"`
}

// CryptoItem holds processed data for a single cryptocurrency
type CryptoItem struct {
	Sym       string
	CoinID    string
	Held      float64
	Entry     float64
	Price     float64
	PrevPrice float64
	OHLCData  []float64
	Index     int // Preserve config order
	RSI       float64
	BBUpper   float64
	BBLower   float64
	Signal    string // BUY, SELL, HOLD
}

// runCrypto is called from main.go via --crypto flag
func RunCrypto(mode string) error {
	// Load config
	config, err := loadCryptoConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get cache paths
	cacheDir := getCacheDir()
	curFile := filepath.Join(cacheDir, "hyprlock-crypto.json")
	prevFile := filepath.Join(cacheDir, "hyprlock-crypto-prev.json")
	ohlcFile := filepath.Join(cacheDir, "hyprlock-ohlc.json")

	// Fetch prices
	ids := make([]string, 0)
	items := make([]CryptoItem, 0)

	for i, p := range config.Pairs {
		if p.Coin == "" {
			continue
		}
		ids = append(ids, p.Coin)
		items = append(items, CryptoItem{
			Sym:    p.Sym,
			CoinID: p.Coin,
			Held:   p.Held,
			Entry:  p.Entry,
			Index:  i,
		})
	}

	// Load current prices first (becomes "prev")
	curPrices := loadJSON(curFile)

	// Save current as prev BEFORE fetching new prices
	prevData := make(map[string]float64)
	for sym, data := range curPrices {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if usd, ok := dataMap["usd"].(float64); ok {
				prevData[sym] = usd
			}
		}
	}
	if len(prevData) > 0 {
		if data, err := json.Marshal(prevData); err == nil {
			os.WriteFile(prevFile, data, 0644)
		}
	}

	// Now fetch fresh prices
	fetchPrices(ids, curFile, config.APIKey)

	// Load fresh prices
	curPrices = loadJSON(curFile)

	// Update items with prices
	for i := range items {
		if curData, ok := curPrices[items[i].CoinID].(map[string]interface{}); ok {
			if usd, ok := curData["usd"].(float64); ok {
				items[i].Price = usd
			}
		}
		// Use the prev data we saved above
		if items[i].CoinID != "" {
			if f, ok := prevData[items[i].CoinID]; ok {
				items[i].PrevPrice = f
			}
		}
	}

	// Load OHLC data for RSI
	loadOHLCData(items, ohlcFile, config.APIKey)

	// Apply default indicator values if not set in config
	if config.Indicators.RSIPeriod == 0 {
		config.Indicators.RSIPeriod = 14
	}
	if config.Indicators.Oversold == 0 {
		config.Indicators.Oversold = 30
	}
	if config.Indicators.Overbought == 0 {
		config.Indicators.Overbought = 70
	}
	if config.Indicators.BBPeriod == 0 {
		config.Indicators.BBPeriod = 20
	}
	if config.Indicators.BBStdDev == 0 {
		config.Indicators.BBStdDev = 2.0
	}

	// Calculate indicators for each item
	for i := range items {
		if len(items[i].OHLCData) > config.Indicators.RSIPeriod {
			items[i].RSI = calculateRSI(items[i].OHLCData, config.Indicators.RSIPeriod)
		}
		if len(items[i].OHLCData) > config.Indicators.BBPeriod {
			items[i].BBUpper, items[i].BBLower = calculateBollingerBands(items[i].OHLCData, config.Indicators.BBPeriod, config.Indicators.BBStdDev)
		}
		items[i].Signal = calculateSignal(items[i].RSI, config.Indicators.Oversold, config.Indicators.Overbought)
	}

	// Route to output mode
	mode = strings.ToUpper(mode)

	switch mode {
	case "ROWML":
		return outputROWML(items, config.Display.ShowTotals, curFile)
	case "ROW":
		return outputROW(items)
	default:
		// Single symbol mode
		for _, item := range items {
			if strings.ToUpper(item.Sym) == mode {
				return outputSingle(item)
			}
		}
		fmt.Println("--")
		return nil
	}
}

func loadCryptoConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(usr.HomeDir, ".config", "crypto.toml")

	// Default config with default pairs - used when config file doesn't exist
	config := &Config{
		Display: DisplaySettings{ShowTotals: false},
		APIKey:  "",
		Pairs: []PairConfig{
			{Sym: "BTC", Coin: "bitcoin", Held: 0, Entry: 0},
			{Sym: "ETH", Coin: "ethereum", Held: 0, Entry: 0},
			{Sym: "XMR", Coin: "monero", Held: 0, Entry: 0},
		},
	}

	_, err = toml.DecodeFile(configPath, config)
	if err != nil {
		// Return default config if file doesn't exist
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	return config, nil
}

func getCacheDir() string {
	usr, _ := user.Current()
	cacheDir := filepath.Join(usr.HomeDir, ".cache")
	os.MkdirAll(cacheDir, 0755)
	return cacheDir
}

func fetchPrices(ids []string, curFile string, apiKey string) {
	if len(ids) == 0 {
		return
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", strings.Join(ids, ","))
	if apiKey != "" {
		url += "&x_cg_demo_api_key=" + apiKey
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "ArchRiot/hyprlock-crypto")

	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if json.NewDecoder(resp.Body).Decode(&data) == nil {
		tmp := curFile + ".tmp"
		f, _ := os.Create(tmp)
		json.NewEncoder(f).Encode(data)
		f.Close()
		os.Rename(tmp, curFile)
	}
}

func loadJSON(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

func loadOHLCData(items []CryptoItem, ohlcFile string, apiKey string) {
	// Check if cache file exists and is valid
	stat, err := os.Stat(ohlcFile)
	cacheValid := false
	cacheAge := 999999

	if err == nil {
		cacheAge = int(time.Since(stat.ModTime()).Seconds())
		if cacheAge < 1800 && stat.Size() > 10 {
			cacheValid = true
		}
	}

	if cacheValid {
		// Load from cache
		raw := loadJSON(ohlcFile)
		for k, v := range raw {
			if arr, ok := v.([]interface{}); ok {
				prices := make([]float64, len(arr))
				for j, p := range arr {
					if f, ok := p.(float64); ok {
						prices[j] = f
					}
				}
				if len(prices) > 0 {
					for i := range items {
						if items[i].Sym == k {
							items[i].OHLCData = prices
						}
					}
				}
			}
		}
	} else {
		// Refresh OHLC data from API
		ohlcCache := make(map[string][]float64)
		for i := range items {
			prices := fetchOHLC(items[i].CoinID, 90, apiKey)
			if len(prices) > 0 {
				ohlcCache[items[i].Sym] = prices
				items[i].OHLCData = prices
			}
		}
		// Save cache only if we got data
		if len(ohlcCache) > 0 {
			data, _ := json.Marshal(ohlcCache)
			os.WriteFile(ohlcFile, data, 0644)
		}
	}
}

func fetchOHLC(coinID string, days int, apiKey string) []float64 {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=usd&days=%d", coinID, days)
	if apiKey != "" {
		url += "&x_cg_demo_api_key=" + apiKey
	}
	// DEBUG: fetchOHLC for %s url=%s")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "ArchRiot/hyprlock-crypto")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// DEBUG: fetchOHLC error for %s: %v")
		return nil
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		// DEBUG: fetchOHLC json error for %s: %v")
		return nil
	}

	if prices, ok := data["prices"].([]interface{}); ok {
		// DEBUG: got %d prices for %s\n", len(prices), coinID)
		result := make([]float64, len(prices))
		for i, p := range prices {
			if arr, ok := p.([]interface{}); ok && len(arr) > 1 {
				if f, ok := arr[1].(float64); ok {
					result[i] = f
				}
			}
		}
		return result
	}
	// DEBUG: no prices in response for %s, data=%v")
	return nil
}

// RSI calculation
func calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 0
	}
	deltas := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		deltas[i-1] = prices[i] - prices[i-1]
	}
	gains := make([]float64, len(deltas))
	losses := make([]float64, len(deltas))
	for i, d := range deltas {
		if d > 0 {
			gains[i] = d
		} else {
			losses[i] = -d
		}
	}
	// Use ALL prior data, not just last 'period' values
	avgGain := 0.0
	avgLoss := 0.0
	for i := range gains {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(len(gains))
	avgLoss /= float64(len(losses))
	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return math.Round(rsi*10) / 10
}

// Calculate signal based on RSI
func calculateSignal(rsi float64, oversold int, overbought int) string {
	if rsi == 0 {
		return "HOLD"
	}
	if rsi < float64(oversold) {
		return "BUY"
	}
	if rsi > float64(overbought) {
		return "SELL"
	}
	return "HOLD"
}

// Calculate Bollinger Bands
func calculateBollingerBands(prices []float64, period int, stdDev float64) (upper float64, lower float64) {
	if len(prices) < period {
		return 0, 0
	}

	recentPrices := prices[len(prices)-period:]

	sum := 0.0
	for _, p := range recentPrices {
		sum += p
	}
	sma := sum / float64(period)

	variance := 0.0
	for _, p := range recentPrices {
		diff := p - sma
		variance += diff * diff
	}
	std := math.Sqrt(variance / float64(period))

	upper = sma + (stdDev * std)
	lower = sma - (stdDev * std)
	return math.Round(upper*100) / 100, math.Round(lower*100) / 100
}

// calculateSellLimit computes the optimal sell limit (max 20% of entry value)
func calculateSellLimit(sym string, currentPrice, entryPrice, held float64, item CryptoItem, items []CryptoItem) string {
	// Calculate max sell value: 20% of entry (entry * held * 0.20)
	maxSellValue := entryPrice * held * 0.20
	if maxSellValue <= 0 || currentPrice <= 0 {
		return ""
	}

	// Calculate max coins we could sell at current price
	maxCoins := maxSellValue / currentPrice

	// Determine best number to sell (20% of holdings, within max)
	var unitsToSell float64
	if held < 1.0 {
		unitsToSell = held * 0.20
		if unitsToSell < 0.01 {
			unitsToSell = held
		}
	} else {
		unitsToSell = math.Min(held*0.20, maxCoins)
		unitsToSell = math.Ceil(unitsToSell*100) / 100
		if unitsToSell < 0.01 {
			unitsToSell = 0.01
		}
	}

	// Calculate target price: entry * 1.20, cap at 90-day high * 0.80
	targetPrice := entryPrice * 1.20
	if len(item.OHLCData) > 0 {
		high := item.OHLCData[0]
		for _, p := range item.OHLCData {
			if p > high {
				high = p
			}
		}
		maxTarget := high * 0.80
		if targetPrice > maxTarget {
				targetPrice = maxTarget
			}
		}
	if targetPrice <= currentPrice {
		targetPrice = currentPrice * 1.20
	}

	unitsStr := formatUnits(unitsToSell)
	targetStr := formatPrice(targetPrice)

	if currentPrice < entryPrice && entryPrice > 0 {
		stopPrice := currentPrice * 0.90
		stopStr := formatPrice(stopPrice)
		return fmt.Sprintf("%s @ $%s (Stop)", unitsStr, stopStr)
	}

	return fmt.Sprintf("%s @ $%s", unitsStr, targetStr)
}

// Output formatters matching shell script exactly
// Shell format: f"{held:>9.2f} x" which gives 11 chars (9 for number + space + x)
// In Python f-string, the width applies to the number only, not the added " x"
// So Go must use: fmt.Sprintf("%9.2f", h) + " x" (not "%9.2f x")
func fmtHeld(h float64) string {
	if h > 0 {
		return fmt.Sprintf("%9.2f", h) + " x"
	}
	return ""
}

func fmtPrice(price float64) string {
	if math.IsNaN(price) || math.IsInf(price, 0) || price == 0 {
		return "$" + strings.Repeat("-", 10)
	}
	// Format: "$" + 10-char number (with comma separators, right-justified)
	// Shell uses: f"${v:>10,.2f}" which is "$" + 10-char number = 11 chars total
	return "$" + formatNumber(price)
}

func fmtPercent(p float64) string {
	if math.IsNaN(p) || math.IsInf(p, 0) {
		return "  ----%"
	}
	// Shell format: f"  {abs(pct):>5.2f}%" for positive (8 chars)
	// Negative: f"- {abs(pct):>5.2f}%" which is "-" + space + 5-char number = 7 chars total
	if p >= 0 {
		return fmt.Sprintf("  %5.2f%%", p)
	}
	// For negative: use abs value with width specifier to match shell
	return fmt.Sprintf("- %5.2f%%", math.Abs(p))
}

func fmtPercentStable() string {
	// For USD stablecoin: 8 chars "   0.00%"
	return "   0.00%"
}

func fmtAmountStable() string {
	// Shell format: f"     --------" = 8 chars with leading spaces, NO trailing space
	return "     --------"
}

func fmtSignedAmount(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "     --------"
	}
	// Format: " " + %12,.2f for positive, "-" + %12,.2f for negative
	if v >= 0 {
		return fmt.Sprintf(" %12.2f", v)
	}
	return fmt.Sprintf("-%12.2f", -v)
}

func fmtSignedAmountWithSpace(v float64) string {
	// Shell format: f" {abs(v):>12,.2f}" = 13 chars (leading space + 12-char number)
	// NO trailing space - matches shell exactly
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "     -------- "
	}
	// Use abs value for formatting to avoid comma being added to negative sign
	numStr := formatNumberWithWidth(math.Abs(v), 12)
	if v >= 0 {
		return " " + numStr
	}
	return "-" + numStr
}

func outputSingle(item CryptoItem) error {
	if item.Price == 0 {
		fmt.Println("--")
		return nil
	}

	entry := item.Entry
	held := item.Held

	if entry > 0 && held > 0 {
		glAmt := (item.Price - entry) * held
		glPct := ((item.Price - entry) / entry) * 100
		fmt.Printf("%s %s %s %s %s\n",
			item.Sym,
			fmtHeld(held),
			fmtPrice(item.Price),
			fmtPercent(glPct),
			fmtSignedAmount(glAmt))
	} else {
		fmt.Printf("%s %s\n", item.Sym, fmtPrice(item.Price))
	}
	return nil
}

// formatNumber formats with 2 decimal places and comma separators, padded to 10 chars
func formatNumber(v float64) string {
	str := fmt.Sprintf("%.2f", v)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}
	var result strings.Builder
	length := len(intPart)
	for i, c := range intPart {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(c)
	}
	intPart = result.String()

	if decPart != "" {
		intPart = intPart + "." + decPart
	}

	// Pad to 10 chars to match shell price formatting
	if len(intPart) < 10 {
		intPart = strings.Repeat(" ", 10-len(intPart)) + intPart
	}
	return intPart
}

// formatNumberWithWidth formats with 2 decimal places and comma separators, padded to specified width
func formatNumberWithWidth(v float64, width int) string {
	str := fmt.Sprintf("%.2f", v)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}
	var result strings.Builder
	length := len(intPart)
	for i, c := range intPart {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(c)
	}
	intPart = result.String()

	if decPart != "" {
		intPart = intPart + "." + decPart
	}

	if len(intPart) < width {
		intPart = strings.Repeat(" ", width-len(intPart)) + intPart
	}
	return intPart
}

func outputROW(items []CryptoItem) error {
	var parts []string
	for _, item := range items {
		if item.Price == 0 {
			continue
		}
		held := item.Held
		var line string
		if held > 0 && item.Entry > 0 {
			glPct := ((item.Price - item.Entry) / item.Entry) * 100
			line = fmt.Sprintf("%s %s %s %s", item.Sym, fmtHeld(held), fmtPrice(item.Price), fmtPercent(glPct))
		} else {
			line = fmt.Sprintf("%s %s", item.Sym, fmtPrice(item.Price))
		}
		parts = append(parts, line)
	}
	fmt.Println(strings.Join(parts, " • "))
	return nil
}

func outputROWML(items []CryptoItem, showTotals bool, curFile string) error {
	lines := []string{}

	// Header row
	header := "COIN     HELD   PRICE        % GAINS       $ GAINS           REBALANCE"
	lines = append(lines, header)
	separator := "------ ------   -----------  -------  ------------   -----------------"
	lines = append(lines, separator)

	// Sort items - preserve config order, move USD to end (matching shell)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Sym == "USD" {
			return false
		}
		if items[j].Sym == "USD" {
			return true
		}
		return items[i].Index < items[j].Index
	})

	// Display up to 6 items (matching shell)
	displayItems := items
	if len(items) > 6 {
		displayItems = items[:6]
	}

	for _, item := range displayItems {
		arrow := " •"
		if item.Price > 0 && item.PrevPrice > 0 {
			if item.Price > item.PrevPrice {
				arrow = " ▲"
			} else if item.Price < item.PrevPrice {
				arrow = " ▼"
			} else {
				arrow = " •"
			}
		}

		// Check for USD stablecoin
		isUSDStable := item.Sym == "USD" && item.Price >= 0.99 && item.Price <= 1.01

		symStr := item.Sym

		var heldStr, pctStr, amtStr string
		if item.Held > 0 && item.Entry > 0 && !isUSDStable {
			glAmt := (item.Price - item.Entry) * item.Held
			glPct := ((item.Price - item.Entry) / item.Entry) * 100
			heldStr = fmtHeld(item.Held)
			pctStr = fmtPercent(glPct)
			amtStr = fmtSignedAmountWithSpace(glAmt)
		} else if isUSDStable {
			heldStr = fmt.Sprintf("%9.2f", item.Held) + " x"
			pctStr = fmtPercentStable()
			amtStr = fmtAmountStable()
		} else {
			heldStr = fmtHeld(item.Held)
			pctStr = ""
			amtStr = ""
		}

		// Calculate sell limit
		sellStr := ""
		if item.Held > 0 {
			sellStr = fmt.Sprintf("%17s", calculateSellLimit(item.Sym, item.Price, item.Entry, item.Held, item, items))
		}

		// Build line matching shell format exactly:
		// Shell: f"{sym} {held_str} {price_str} {pct_str} {amt_str}{arrow}{sell_str:>22}"
		// Note: NO space between amt_str and arrow in shell
		priceStr := fmtPrice(item.Price)
		line := fmt.Sprintf("%s %s %s %s %s%s %s", symStr, heldStr, priceStr, pctStr, amtStr, arrow, sellStr)
		lines = append(lines, line)
	}

	// Calculate totals
	var heldTotal, gainTotal float64
	haveValue := false
	haveGain := false
	for _, item := range items {
		if item.Price == 0 {
			continue
		}
		if item.Held > 0 {
			heldTotal += item.Held * item.Price
			haveValue = true
		}
		if item.Held > 0 && item.Entry > 0 {
			gainTotal += (item.Price - item.Entry) * item.Held
			haveGain = true
		}
	}

	if haveGain {
		if showTotals && haveValue {
			heldStr := "$ " + formatNumberWithWidth(heldTotal, 10)
			var gainStr string
			if gainTotal >= 0 {
				gainStr = " " + formatNumberWithWidth(gainTotal, 12)
			} else {
				gainStr = "-" + formatNumberWithWidth(-gainTotal, 12)
			}
			lines = append(lines, fmt.Sprintf("%s%s%s", strings.Repeat(" ", 37), gainStr, strings.Repeat(" ", 15))+heldStr)
		} else {
			var gainStr string
			if gainTotal >= 0 {
				gainStr = " " + formatNumberWithWidth(gainTotal, 12)
			} else {
				gainStr = "-" + formatNumberWithWidth(-gainTotal, 12)
			}
			lines = append(lines, fmt.Sprintf("%s%s%s", strings.Repeat(" ", 37), gainStr, strings.Repeat(" ", 23)))
		}
	}

	fmt.Println(strings.Join(lines, "\n"))

	// Save prev snapshot
	prevFile := strings.Replace(curFile, "crypto.json", "crypto-prev.json", 1)
	prevData := make(map[string]float64)
	for _, item := range items {
		if item.Price > 0 {
			prevData[item.Sym] = item.Price
		}
	}
	data, _ := json.Marshal(prevData)
	os.WriteFile(prevFile, data, 0644)

	return nil
}
