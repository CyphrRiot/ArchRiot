package waybar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EmitTemp will emit temperature status JSON (moved from --waybar-temp in main.go).
func EmitTemp(args []string) int {
	findSensor := func() string {
		if entries, err := os.ReadDir("/sys/class/hwmon"); err == nil {
			for _, e := range entries {
				hp := filepath.Join("/sys/class/hwmon", e.Name())
				nb, err := os.ReadFile(filepath.Join(hp, "name"))
				if err != nil {
					continue
				}
				name := strings.TrimSpace(string(nb))
				switch name {
				case "coretemp", "k10temp", "zenpower":
					tf := filepath.Join(hp, "temp1_input")
					if st, err := os.Stat(tf); err == nil && !st.IsDir() {
						return tf
					}
				}
			}
		}
		if entries, err := os.ReadDir("/sys/class/thermal"); err == nil {
			for _, e := range entries {
				if !strings.HasPrefix(e.Name(), "thermal_zone") {
					continue
				}
				zp := filepath.Join("/sys/class/thermal", e.Name())
				tb, err := os.ReadFile(filepath.Join(zp, "type"))
				if err != nil {
					continue
				}
				if strings.TrimSpace(string(tb)) == "x86_pkg_temp" {
					tf := filepath.Join(zp, "temp")
					if st, err := os.Stat(tf); err == nil && !st.IsDir() {
						return tf
					}
				}
			}
		}
		tf := "/sys/class/thermal/thermal_zone0/temp"
		if st, err := os.Stat(tf); err == nil && !st.IsDir() {
			return tf
		}
		return ""
	}
	readTempC := func() (float64, bool) {
		s := findSensor()
		if s == "" {
			return 0, false
		}
		b, err := os.ReadFile(s)
		if err != nil {
			return 0, false
		}
		raw := strings.TrimSpace(string(b))
		val, _ := strconv.ParseFloat(raw, 64)
		if val == 0 {
			if iv, e := strconv.ParseInt(raw, 10, 64); e == nil {
				if iv > 1000 {
					return float64(iv) / 1000.0, true
				}
				return float64(iv), true
			}
			return 0, false
		}
		if val > 1000 {
			return val / 1000.0, true
		}
		return val, true
	}
	tempC, ok := readTempC()
	if !ok {
		fmt.Println(`{"text":"-- 󰈸","tooltip":"Temperature sensor not available","class":"critical","percentage":0}`)
		return 0
	}
	tempPct := (tempC - 60.0) * 100.0 / 35.0
	if tempPct < 0 {
		tempPct = 0
	}
	if tempPct > 100 {
		tempPct = 100
	}
	bar := func(p float64) string {
		switch {
		case p <= 10:
			return "▁"
		case p <= 25:
			return "▂"
		case p <= 40:
			return "▃"
		case p <= 55:
			return "▄"
		case p <= 70:
			return "▅"
		case p <= 80:
			return "▆"
		case p <= 90:
			return "▇"
		default:
			return "█"
		}
	}(tempPct)
	class := "normal"
	switch {
	case tempC >= 90:
		class = "critical"
	case tempC >= 80:
		class = "warning"
	}
	out := Output{
		Text:       fmt.Sprintf("%s 󰈸", bar),
		Tooltip:    fmt.Sprintf("CPU Temperature: %.1f°C", tempC),
		Class:      class,
		Percentage: int(tempC + 0.5),
	}
	if js, err := json.Marshal(out); err == nil {
		fmt.Println(string(js))
	} else {
		fmt.Println(`{"text":"-- 󰈸","tooltip":"Temperature Error: marshal","class":"critical","percentage":0}`)
	}
	return 0
}
