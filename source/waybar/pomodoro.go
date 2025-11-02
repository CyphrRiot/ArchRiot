package waybar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
