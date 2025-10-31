package waybar

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Package waybar provides JSON emitters for Waybar modules.
// This file intentionally contains only scaffolding to enable
// incremental extraction from main.go without any behavior change yet.
//
// Contract for emitters:
//   - Each function writes the appropriate JSON to stdout (once) when implemented.
//   - Each function returns an exit code (0 = success), but must not call os.Exit.
//   - Arguments are passed through for future flag handling when needed.

type Output struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip"`
	Class      string `json:"class"`
	Percentage int    `json:"percentage"`
}

// EmitPomodoro will emit Pomodoro status JSON (moved from --waybar-pomodoro in main.go).
func EmitPomodoro(args []string) int {
	type pomoState struct {
		Mode            string   `json:"mode"` // "idle", "work", "break", "break_complete"
		Running         bool     `json:"running"`
		EndTimeISO      string   `json:"end_time"`
		PausedRemaining *float64 `json:"paused_remaining"`
	}
	type pomoOut struct {
		Text       string `json:"text"`
		Tooltip    string `json:"tooltip"`
		Class      string `json:"class"`
		Percentage int    `json:"percentage"`
	}

	now := time.Now()
	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".config", "archriot", "archriot.conf")
	timerFile := "/tmp/waybar-tomato.json"
	stateFile := "/tmp/waybar-tomato-timer.state"

	enabled := true
	workMinutes := 25
	breakMinutes := 5

	// Minimal INI parse for [pomodoro]
	if b, err := os.ReadFile(configPath); err == nil {
		inSection := false
		for _, l := range strings.Split(string(b), "\n") {
			line := strings.TrimSpace(l)
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				continue
			}
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				sect := strings.ToLower(strings.Trim(line, "[]"))
				inSection = sect == "pomodoro"
				continue
			}
			if !inSection {
				continue
			}
			if strings.HasPrefix(strings.ToLower(line), "enabled") && strings.Contains(line, "=") {
				val := strings.ToLower(strings.TrimSpace(strings.SplitN(line, "=", 2)[1]))
				if val == "false" || val == "0" || val == "no" || val == "off" {
					enabled = false
				}
			}
			if strings.HasPrefix(strings.ToLower(line), "duration") && strings.Contains(line, "=") {
				val := strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
				if n, err := strconv.Atoi(val); err == nil && n > 0 && n <= 120 {
					workMinutes = n
				}
			}
		}
	}

	// Helper: JSON marshal or fallback
	mustJSON := func(v interface{}) []byte {
		j, err := json.Marshal(v)
		if err != nil {
			return []byte("{}")
		}
		return j
	}

	// Load current timer state
	state := pomoState{Mode: "idle", Running: false, EndTimeISO: "", PausedRemaining: nil}
	if b, err := os.ReadFile(timerFile); err == nil {
		_ = json.Unmarshal(b, &state)
	}

	// Helper to persist state
	save := func() {
		_ = os.WriteFile(timerFile, mustJSON(state), 0644)
	}

	// Process click action if present
	if _, err := os.Stat(stateFile); err == nil {
		// Read and remove click file
		var act struct {
			Action string `json:"action"`
		}
		if b, err := os.ReadFile(stateFile); err == nil {
			_ = json.Unmarshal(b, &act)
		}
		_ = os.Remove(stateFile)

		if enabled {
			switch act.Action {
			case "reset":
				state.Mode = "idle"
				state.Running = false
				state.EndTimeISO = ""
				state.PausedRemaining = nil
				save()
			case "toggle":
				switch {
				case state.Mode == "idle" || state.Mode == "break_complete":
					state.Mode = "work"
					state.Running = true
					state.EndTimeISO = now.Add(time.Duration(workMinutes) * time.Minute).Format(time.RFC3339)
					state.PausedRemaining = nil
					save()
				case state.Running:
					// Pause
					if state.EndTimeISO != "" {
						if t, err := time.Parse(time.RFC3339, state.EndTimeISO); err == nil {
							rem := t.Sub(now).Seconds()
							if rem < 0 {
								rem = 0
							}
							r := rem
							state.PausedRemaining = &r
						}
					}
					state.Running = false
					save()
				default:
					// Resume
					if state.PausedRemaining != nil && *state.PausedRemaining > 0 {
						end := now.Add(time.Duration(*state.PausedRemaining) * time.Second)
						state.EndTimeISO = end.Format(time.RFC3339)
						state.PausedRemaining = nil
					}
					state.Running = true
					save()
				}
			}
		}
	}

	// Compute remaining and handle transitions
	remaining := 0.0
	total := 0.0
	switch state.Mode {
	case "work":
		total = float64(workMinutes * 60)
	case "break":
		total = float64(breakMinutes * 60)
	}

	if state.Running && state.EndTimeISO != "" {
		if t, err := time.Parse(time.RFC3339, state.EndTimeISO); err == nil {
			remaining = t.Sub(now).Seconds()
		}
	} else if !state.Running && state.PausedRemaining != nil {
		remaining = *state.PausedRemaining
	}

	if remaining <= 0 {
		switch state.Mode {
		case "work":
			// Transition to break
			state.Mode = "break"
			state.Running = true
			state.EndTimeISO = now.Add(time.Duration(breakMinutes) * time.Minute).Format(time.RFC3339)
			state.PausedRemaining = nil
			save()
			remaining = float64(breakMinutes * 60)
			total = remaining
		case "break":
			// Break complete
			state.Mode = "break_complete"
			state.Running = false
			state.EndTimeISO = ""
			state.PausedRemaining = nil
			save()
		default:
			remaining = 0
		}
	}

	// Build output
	out := pomoOut{Text: "", Tooltip: "", Class: "idle", Percentage: 0}
	if !enabled {
		out.Text = "󰌾 --:--"
		out.Tooltip = "Pomodoro Timer - Disabled"
		out.Class = "disabled"
		fmt.Println(string(mustJSON(out)))
		return 0
	}

	switch state.Mode {
	case "idle":
		out.Text = fmt.Sprintf("󰌾 %02d:00", workMinutes)
		out.Tooltip = "Pomodoro Timer - Click to start"
		out.Class = "idle"
	case "break_complete":
		out.Text = "󰌾 Ready"
		out.Tooltip = "Break over! Click to start next session"
		out.Class = "break_complete"
	case "work", "break":
		rem := int(remaining + 0.5)
		if rem < 0 {
			rem = 0
		}
		mins := rem / 60
		secs := rem % 60
		icon := "󰔛"
		if state.Mode == "break" {
			icon = "☕"
		}
		if !state.Running {
			icon = "󰏤"
		}
		status := "Work"
		if state.Mode == "break" {
			status = "Break"
		}
		if !state.Running {
			status = "Paused"
		}
		out.Text = fmt.Sprintf("%s %02d:%02d", icon, mins, secs)
		out.Tooltip = fmt.Sprintf("%s - %02d:%02d remaining", status, mins, secs)
		if !state.Running {
			out.Class = "paused"
		} else {
			out.Class = state.Mode
		}
		if total > 0 {
			progress := 100.0 - ((float64(rem) / total) * 100.0)
			if progress < 0 {
				progress = 0
			}
			if progress > 100 {
				progress = 100
			}
			out.Percentage = int(progress + 0.5)
		}
	}

	fmt.Println(string(mustJSON(out)))
	return 0
}

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

// EmitCPU will emit CPU usage JSON (moved from --waybar-cpu in main.go).
func EmitCPU(args []string) int {
	type cpuOut struct {
		Text       string `json:"text"`
		Tooltip    string `json:"tooltip"`
		Class      string `json:"class"`
		Percentage int    `json:"percentage"`
	}
	readCPU := func() (idle, total uint64, ok bool) {
		data, err := os.ReadFile("/proc/stat")
		if err != nil {
			return 0, 0, false
		}
		for _, ln := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(ln, "cpu ") {
				fields := strings.Fields(ln)
				if len(fields) < 8 {
					return 0, 0, false
				}
				parse := func(s string) uint64 {
					v, _ := strconv.ParseUint(s, 10, 64)
					return v
				}
				user := parse(fields[1])
				nice := parse(fields[2])
				system := parse(fields[3])
				idleVal := parse(fields[4])
				iowait := parse(fields[5])
				irq := parse(fields[6])
				softirq := parse(fields[7])
				steal := uint64(0)
				if len(fields) > 8 {
					steal = parse(fields[8])
				}
				idle := idleVal + iowait
				total := user + nice + system + idleVal + iowait + irq + softirq + steal
				return idle, total, true
			}
		}
		return 0, 0, false
	}
	i1, t1, ok := readCPU()
	if !ok {
		fmt.Println(`{"text":"-- 󰍛","tooltip":"CPU Error: cannot read /proc/stat","class":"critical","percentage":0}`)
		return 0
	}
	time.Sleep(120 * time.Millisecond)
	i2, t2, ok := readCPU()
	if !ok || t2 <= t1 || i2 < i1 {
		fmt.Println(`{"text":"-- 󰍛","tooltip":"CPU Error: invalid delta","class":"critical","percentage":0}`)
		return 0
	}
	dIdle := i2 - i1
	dTotal := t2 - t1
	usage := (float64(dTotal-dIdle) / float64(dTotal)) * 100.0
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
	}(usage)
	class := "normal"
	if usage >= 90 {
		class = "critical"
	} else if usage >= 75 {
		class = "warning"
	}
	out := Output{
		Text:       fmt.Sprintf("%s 󰍛", bar),
		Tooltip:    fmt.Sprintf("CPU Usage: %.1f%%", usage),
		Class:      class,
		Percentage: int(usage + 0.5),
	}
	if js, err := json.Marshal(out); err == nil {
		fmt.Println(string(js))
	} else {
		fmt.Println(`{"text":"-- 󰍛","tooltip":"CPU Error: marshal","class":"critical","percentage":0}`)
	}
	return 0
}

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

// EmitVolume will emit speaker volume JSON (moved from --waybar-volume in main.go).
func EmitVolume(args []string) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	vol := func() (int, bool) {
		if have("wpctl") {
			out, err := exec.Command("sh", "-lc", "wpctl get-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | awk '{print int($2*100+0.5)}'").Output()
			if err == nil {
				v, _ := strconv.Atoi(strings.TrimSpace(string(out)))
				return v, true
			}
		}
		if have("pamixer") {
			// Ensure PipeWire/Pulse is ready
			if have("pactl") && exec.Command("pactl", "info").Run() != nil {
				return 0, false
			}
			out, err := exec.Command("pamixer", "--get-volume").Output()
			if err == nil {
				v, _ := strconv.Atoi(strings.TrimSpace(string(out)))
				return v, true
			}
		}
		if have("pactl") {
			out, err := exec.Command("sh", "-lc", "pactl get-sink-volume @DEFAULT_SINK@ | grep -o '[0-9]\\+%' | head -n1 | tr -d '%'").Output()
			if err == nil {
				v, _ := strconv.Atoi(strings.TrimSpace(string(out)))
				return v, true
			}
		}
		return 0, false
	}

	muted := func() (bool, bool) {
		if have("wpctl") {
			err := exec.Command("sh", "-lc", "wpctl get-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | grep -qi yes").Run()
			return err == nil, true
		}
		if have("pamixer") {
			// Ensure PipeWire/Pulse is ready
			if have("pactl") && exec.Command("pactl", "info").Run() != nil {
				return false, false
			}
			out, err := exec.Command("pamixer", "--get-mute").Output()
			if err == nil {
				return strings.TrimSpace(string(out)) == "true", true
			}
		}
		if have("pactl") {
			err := exec.Command("sh", "-lc", "pactl get-sink-mute @DEFAULT_SINK@ | grep -qi yes").Run()
			return err == nil, true
		}
		return false, false
	}

	v, vok := vol()
	m, mok := muted()
	if !vok || !mok {
		fmt.Println(`{"text":"▁ 󰖁","tooltip":"Audio not ready","class":"muted","percentage":0}`)
		return 0
	}
	if m {
		fmt.Println(`{"text":"▁ 󰖁","tooltip":"Speaker: Muted","class":"muted","percentage":0}`)
		return 0
	}

	var bar string
	switch {
	case v <= 2:
		bar = "▁"
	case v <= 5:
		bar = "▂"
	case v <= 10:
		bar = "▃"
	case v <= 20:
		bar = "▄"
	case v <= 35:
		bar = "▅"
	case v <= 50:
		bar = "▆"
	case v <= 75:
		bar = "▇"
	default:
		bar = "█"
	}

	icon := "󰕾"
	if v == 0 {
		icon = "󰕿"
	} else if v <= 33 {
		icon = "󰖀"
	} else if v <= 66 {
		icon = "󰕾"
	}

	class := "normal"
	if v >= 100 {
		class = "critical"
	} else if v >= 85 {
		class = "warning"
	}

	fmt.Println(fmt.Sprintf(`{"text":"%s %s","tooltip":"Speaker Volume: %d%%","class":"%s","percentage":%d}`, bar, icon, v, class, v))
	return 0
}
