package waybar

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EmitMemory will emit memory usage JSON (moved from --waybar-memory in main.go).
func EmitMemory(args []string) int {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println(`{"text":"-- 󰾆","tooltip":"Memory Error: cannot read /proc/meminfo","class":"critical","percentage":0}`)
		return 0
	}
	parse := func(key string) int64 {
		for _, ln := range strings.Split(string(data), "\n") {
			fields := strings.Fields(ln)
			if len(fields) >= 2 && strings.TrimSuffix(fields[0], ":") == key {
				val, _ := strconv.ParseInt(fields[1], 10, 64)
				return val
			}
		}
		return 0
	}
	memTotal := parse("MemTotal")
	memAvailable := parse("MemAvailable")
	memFree := parse("MemFree")
	buffers := parse("Buffers")
	cached := parse("Cached")

	if memTotal <= 0 {
		fmt.Println(`{"text":"-- 󰾆","tooltip":"Memory Error: invalid totals","class":"critical","percentage":0}`)
		return 0
	}

	usedModernKB := memTotal - memAvailable
	usedTraditionalKB := memTotal - memFree - buffers - cached

	totalGB := float64(memTotal) / (1024.0 * 1024.0)
	availGB := float64(memAvailable) / (1024.0 * 1024.0)
	usedModernGB := float64(usedModernKB) / (1024.0 * 1024.0)
	usedTraditionalGB := float64(usedTraditionalKB) / (1024.0 * 1024.0)

	percent := (float64(usedTraditionalKB) / float64(memTotal)) * 100.0

	bar := func(p float64) string {
		switch {
		case p <= 0:
			return ""
		case p <= 15:
			return "▁"
		case p <= 30:
			return "▂"
		case p <= 45:
			return "▃"
		case p <= 60:
			return "▄"
		case p <= 75:
			return "▅"
		case p <= 85:
			return "▆"
		case p <= 95:
			return "▇"
		default:
			return "█"
		}
	}(percent)

	class := "normal"
	if percent >= 90 {
		class = "critical"
	} else if percent >= 75 {
		class = "warning"
	}

	out := Output{
		Text:       fmt.Sprintf("%s 󰾆", bar),
		Tooltip:    fmt.Sprintf("Used (Modern): %.1fGB\nUsed (Traditional): %.1fGB\nAvailable: %.1fGB\nTotal: %.1fGB (%.1f%%)", usedModernGB, usedTraditionalGB, availGB, totalGB, percent),
		Class:      class,
		Percentage: int(percent + 0.5),
	}
	if js, err := json.Marshal(out); err == nil {
		fmt.Println(string(js))
	} else {
		fmt.Println(`{"text":"-- 󰾆","tooltip":"Memory Error: marshal","class":"critical","percentage":0}`)
	}
	return 0
}
