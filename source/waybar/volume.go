package waybar

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

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
