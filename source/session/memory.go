package session

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// MemoryCleanOpts controls cache-drop behavior.
type MemoryCleanOpts struct {
	// IfLow runs the clean only when MemAvailable < ThresholdMB.
	IfLow bool
	// ThresholdMB is the MemAvailable threshold in MB for IfLow (default ~1024 when <= 0).
	ThresholdMB int
	// Notify shows a desktop notification (best-effort) when true.
	Notify bool
}

// MemoryClean drops Linux page caches similar to the shell snippet:
//
//	sync
//	echo 3 > /proc/sys/vm/drop_caches
//
// Behavior:
// - Prints memory before/after (free -h if available; /proc fallback otherwise).
// - Uses root direct write when EUID==0.
// - Else tries pkexec (GUI prompt), then sudo.
// - Optional low-memory gate (IfLow/ThresholdMB) to avoid unnecessary cache drops.
// - Optional notify-send on success/failure (best-effort).
//
// Returns 0 on success/bypass; 1 when elevation failed or /proc write was not possible.
func MemoryClean(opts MemoryCleanOpts) int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
	notify := func(msg string) {
		if opts.Notify && have("notify-send") {
			_ = exec.Command("notify-send", "-t", "2000", "ArchRiot", msg).Start()
		}
	}

	// Detect MemAvailable (kB)
	memAvailKB := readMemAvailableKB()

	// Low-memory guard (defaults to 1024 MB if threshold not set)
	if opts.IfLow {
		thrKB := 1024 * 1024
		if opts.ThresholdMB > 0 {
			thrKB = opts.ThresholdMB * 1024
		}
		if memAvailKB >= thrKB {
			notify("Memory clean skipped (above threshold)")
			return 0
		}
	}

	printFree("Memory before clearing cache:")

	// Flush FS buffers
	_ = exec.Command("sync").Run()

	// Drop caches with best-available elevation
	if err := dropCaches(); err != nil {
		notify("Memory clean failed (elevation required)")
		fmt.Fprintln(os.Stderr, "memory-clean: failed to drop caches:", err)
		return 1
	}

	printFree("Memory after clearing cache:")
	notify("Memory clean completed")
	return 0
}

func printFree(header string) {
	fmt.Println(header)
	if _, err := exec.LookPath("free"); err == nil {
		cmd := exec.Command("free", "-h")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	} else {
		// Minimal fallback using /proc/meminfo
		mt, ma := readMemTotalKB(), readMemAvailableKB()
		fmt.Printf("MemTotal:     %.2f GiB\n", float64(mt)/1048576.0)
		fmt.Printf("MemAvailable: %.2f GiB\n", float64(ma)/1048576.0)
	}
	fmt.Println()
}

func readMemAvailableKB() int { return readMeminfoKB("MemAvailable:") }
func readMemTotalKB() int     { return readMeminfoKB("MemTotal:") }

func readMeminfoKB(key string) int {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, key) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if v, err := strconv.Atoi(fields[1]); err == nil {
					return v // already in kB
				}
			}
		}
	}
	return 0
}

// dropCaches writes "3" to /proc/sys/vm/drop_caches with best-available elevation.
// Order: direct (root) -> pkexec (GUI) -> sudo (TTY).
func dropCaches() error {
	// If already root, write directly
	if os.Geteuid() == 0 {
		// Writing permissions on /proc are special; the perm argument is ignored.
		return os.WriteFile("/proc/sys/vm/drop_caches", []byte("3\n"), 0o644)
	}

	// Prefer pkexec for GUI prompt (Wayland-friendly; no TTY needed)
	if _, err := exec.LookPath("pkexec"); err == nil {
		cmd := exec.Command("pkexec", "sh", "-lc", "sync; echo 3 > /proc/sys/vm/drop_caches")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("pkexec failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
		}
		return nil
	}

	// Fallback to sudo for terminal environments
	if _, err := exec.LookPath("sudo"); err == nil {
		cmd := exec.Command("sudo", "sh", "-lc", "sync; echo 3 > /proc/sys/vm/drop_caches")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("sudo failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
		}
		return nil
	}

	return fmt.Errorf("no elevation helper available (need root, pkexec, or sudo)")
}
