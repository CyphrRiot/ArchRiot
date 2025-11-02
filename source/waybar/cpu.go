package waybar

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
