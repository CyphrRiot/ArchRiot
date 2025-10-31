package audio

// Package audio provides speaker and microphone volume controls extracted from main.go.
//
// Usage (delegated from CLI):
//   archriot --volume toggle|inc|dec|get
//   archriot --volume mic-toggle|mic-inc|mic-dec|mic-get
//
// Behavior:
// - Prefers PipeWire via wpctl, then pamixer, then pactl.
// - Shows desktop notifications when possible (notify-send).
// - Never calls os.Exit; returns an exit code for the caller to handle.

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run executes the volume subcommands and returns an exit code (0 = success).
func Run(args []string) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }

	notify := func(title, msg, icon string) {
		if have("makoctl") {
			_ = exec.Command("makoctl", "dismiss", "--all").Run()
		}
		if have("notify-send") {
			_ = exec.Command("notify-send",
				"--replace-id=8888",
				"--app-name=Volume Control",
				"--urgency=normal",
				"--icon", icon,
				title, msg,
			).Start()
		}
	}

	useWpctl := have("wpctl")
	usePamixer := have("pamixer")
	usePactl := have("pactl")

	if !useWpctl && !usePamixer && !usePactl {
		notify("Volume Error", "pamixer/wpctl/pactl not found", "dialog-error")
		fmt.Fprintln(os.Stderr, "Error: pamixer/wpctl/pactl is not installed")
		return 1
	}

	vol := func() string {
		// Prefer PipeWire (wpctl) first, then pamixer, then pactl
		if useWpctl {
			out, err := exec.Command("sh", "-lc", "wpctl get-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | awk '{print int($2*100+0.5)}'").Output()
			if err != nil {
				return "0"
			}
			return strings.TrimSpace(string(out))
		}
		if usePamixer {
			out, err := exec.Command("pamixer", "--get-volume").Output()
			if err != nil {
				return "0"
			}
			return strings.TrimSpace(string(out))
		}
		// pactl fallback
		out, err := exec.Command("sh", "-lc", "pactl get-sink-volume @DEFAULT_SINK@ | grep -o '[0-9]\\+%' | head -n1 | tr -d '%'").Output()
		if err != nil {
			return "0"
		}
		return strings.TrimSpace(string(out))
	}

	isMuted := func() bool {
		// Prefer wpctl (PipeWire), then pamixer, then pactl
		if useWpctl {
			return exec.Command("sh", "-lc", "wpctl get-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | grep -qi yes").Run() == nil
		}
		if usePamixer {
			return exec.Command("sh", "-lc", "pamixer --get-mute | grep -q true").Run() == nil
		}
		return exec.Command("sh", "-lc", "pactl get-sink-mute @DEFAULT_SINK@ | grep -qi yes").Run() == nil
	}

	micVol := func() string {
		// Prefer wpctl first, then pamixer, then pactl
		if useWpctl {
			out, err := exec.Command("sh", "-lc", "wpctl get-volume $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | awk '{print int($2*100+0.5)}'").Output()
			if err != nil {
				return "0"
			}
			return strings.TrimSpace(string(out))
		}
		if usePamixer {
			out, err := exec.Command("pamixer", "--default-source", "--get-volume").Output()
			if err != nil {
				return "0"
			}
			return strings.TrimSpace(string(out))
		}
		// pactl fallback
		out, err := exec.Command("sh", "-lc", "pactl get-source-volume @DEFAULT_SOURCE@ | grep -o '[0-9]\\+%' | head -n1 | tr -d '%'").Output()
		if err != nil {
			return "0"
		}
		return strings.TrimSpace(string(out))
	}

	micMuted := func() bool {
		// Prefer wpctl (PipeWire), then pamixer, then pactl
		if useWpctl {
			return exec.Command("sh", "-lc", "wpctl get-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") | grep -qi yes").Run() == nil
		}
		if usePamixer {
			return exec.Command("sh", "-lc", "pamixer --default-source --get-mute | grep -q true").Run() == nil
		}
		return exec.Command("sh", "-lc", "pactl get-source-mute @DEFAULT_SOURCE@ | grep -qi yes").Run() == nil
	}

	usage := func() int {
		fmt.Fprintln(os.Stderr, "Usage: archriot --volume [toggle|inc|dec|get|mic-toggle|mic-inc|mic-dec|mic-get]")
		return 1
	}

	if len(args) < 1 {
		return usage()
	}

	switch args[0] {
	case "toggle":
		// Prefer PipeWire native control first
		if useWpctl {
			_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") toggle").Run()
		} else if usePamixer {
			_ = exec.Command("pamixer", "--toggle-mute").Run()
		} else {
			_ = exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle").Run()
		}
		if isMuted() {
			notify("Audio Muted", "Speaker: Muted", "audio-volume-muted")
		} else {
			notify("Audio Unmuted", "Speaker: "+vol()+"%", "audio-volume-high")
		}
		return 0

	case "inc":
		if useWpctl {
			_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
			_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%+").Run()
		} else if usePamixer {
			_ = exec.Command("pamixer", "--increase", "5", "--unmute").Run()
		} else {
			_ = exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "0").Run()
			_ = exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "+5%").Run()
		}
		notify("Volume Up", "Speaker: "+vol()+"%", "audio-volume-high")
		return 0

	case "dec":
		if useWpctl {
			_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
			_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sinks:/,/Sources:/{if (/Sources:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%-").Run()
		} else if usePamixer {
			_ = exec.Command("pamixer", "--decrease", "5", "--unmute").Run()
		} else {
			_ = exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "0").Run()
			_ = exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", "-5%").Run()
		}
		notify("Volume Down", "Speaker: "+vol()+"%", "audio-volume-low")
		return 0

	case "get":
		fmt.Println(vol())
		return 0

	case "mic-toggle":
		if useWpctl {
			_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") toggle").Run()
		} else if usePamixer {
			_ = exec.Command("pamixer", "--default-source", "--toggle-mute").Run()
		} else {
			_ = exec.Command("pactl", "set-source-mute", "@DEFAULT_SOURCE@", "toggle").Run()
		}
		if micMuted() {
			notify("Microphone Muted", "Microphone: Muted", "microphone-sensitivity-muted")
		} else {
			notify("Microphone Unmuted", "Microphone: "+micVol()+"%", "microphone-sensitivity-high")
		}
		return 0

	case "mic-inc":
		if useWpctl {
			_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
			_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%+").Run()
		} else if usePamixer {
			_ = exec.Command("pamixer", "--default-source", "--increase", "5", "--unmute").Run()
		} else {
			_ = exec.Command("pactl", "set-source-mute", "@DEFAULT_SOURCE@", "0").Run()
			_ = exec.Command("pactl", "set-source-volume", "@DEFAULT_SOURCE@", "+5%").Run()
		}
		notify("Microphone Volume", "Microphone: "+micVol()+"%", "microphone-sensitivity-high")
		return 0

	case "mic-dec":
		if useWpctl {
			_ = exec.Command("sh", "-lc", "wpctl set-mute $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 0").Run()
			_ = exec.Command("sh", "-lc", "wpctl set-volume $(wpctl status | awk '/Sources:/,/Filters:/{if (/Filters:/) exit; print}' | grep '\\*' | awk '{print $3}' | tr -d \".\") 5%-").Run()
		} else if usePamixer {
			_ = exec.Command("pamixer", "--default-source", "--decrease", "5", "--unmute").Run()
		} else {
			_ = exec.Command("pactl", "set-source-mute", "@DEFAULT_SOURCE@", "0").Run()
			_ = exec.Command("pactl", "set-source-volume", "@DEFAULT_SOURCE@", "-5%").Run()
		}
		notify("Microphone Volume", "Microphone: "+micVol()+"%", "microphone-sensitivity-low")
		return 0

	case "mic-get":
		fmt.Println(micVol())
		return 0

	default:
		return usage()
	}
}
