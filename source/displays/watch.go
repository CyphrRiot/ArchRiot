package displays

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// WatchHotplug monitors Hyprland's monitor topology and enforces the external-preferred
// display policy when a topology change is detected. It never triggers a compositor reload.
// Intended to be run as a long-lived background process (systemd --user recommended).
//
// Behavior:
// - Polls `hyprctl -j monitors` at a configurable interval
// - Computes a signature of the raw JSON; on change, waits for a short debounce and runs Enforce()
// - Best-effort; failures to read topology will be retried on the next tick
//
// Flags (pass via args):
//
//	--interval-ms  N   Polling interval in milliseconds (default 800)
//	--debounce-ms  N   Debounce after a change before enforcing (default 500)
//	--verbose          Print change/enforcement logs to stdout
//	--once             Enforce once if hyprctl is ready and exit (no watching)
//	--timeout-ms  N    Only valid with --once. Fail if no hyprctl ready by timeout (default 3000)
//
// Returns 0 on normal termination; nonzero on immediate startup errors.
func WatchHotplug(args []string) int {
	var (
		intervalMS = 800
		debounceMS = 500
		verbose    = false
		once       = false
		timeoutMS  = 3000
	)

	// Manual flag parsing to keep dependencies minimal
	fs := flag.NewFlagSet("watch-hotplug", flag.ContinueOnError)
	fs.SetOutput(nil) // silence default usage output
	fs.IntVar(&intervalMS, "interval-ms", intervalMS, "")
	fs.IntVar(&debounceMS, "debounce-ms", debounceMS, "")
	fs.BoolVar(&verbose, "verbose", verbose, "")
	fs.BoolVar(&once, "once", once, "")
	fs.IntVar(&timeoutMS, "timeout-ms", timeoutMS, "")

	_ = fs.Parse(args)

	// Sanitize values
	if intervalMS < 100 {
		intervalMS = 100
	}
	if debounceMS < 100 {
		debounceMS = 100
	}
	if timeoutMS < 0 {
		timeoutMS = 0
	}

	// Ensure hyprctl is available
	if _, err := exec.LookPath("hyprctl"); err != nil {
		fmt.Fprintln(os.Stderr, "watch-hotplug: hyprctl not found in PATH")
		return 1
	}

	// Helper: wait until hyprctl -j monitors responds or timeout
	waitHypr := func(timeout time.Duration) bool {
		deadline := time.Now().Add(timeout)
		for {
			if sig := readTopologySignature(); sig != "" {
				return true
			}
			if timeout > 0 && time.Now().After(deadline) {
				return false
			}
			time.Sleep(200 * time.Millisecond)
		}
	}

	// For --once, just enforce once when ready and exit
	if once {
		if timeoutMS == 0 || waitHypr(time.Duration(timeoutMS)*time.Millisecond) {
			if verbose {
				fmt.Println("watch-hotplug: hyprctl ready; enforcing once")
			}
			_ = Enforce()
			return 0
		}
		fmt.Fprintln(os.Stderr, "watch-hotplug: hyprctl not ready before timeout ("+strconv.Itoa(timeoutMS)+"ms)")
		return 2
	}

	// Long-running watcher: set up signal handling for graceful exit
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for hyprctl readiness (no hard fail; just log)
	if !waitHypr(5 * time.Second) {
		if verbose {
			fmt.Println("watch-hotplug: hyprctl not ready after 5s; will retry in background")
		}
	}

	lastSig := readTopologySignature()
	if verbose && lastSig != "" {
		fmt.Println("watch-hotplug: initial topology signature:", lastSig[:min(8, len(lastSig))])
	}

	ticker := time.NewTicker(time.Duration(intervalMS) * time.Millisecond)
	defer ticker.Stop()

	debouncing := false
	debounceTimer := time.NewTimer(time.Hour) // placeholder; will be reset
	_ = debounceTimer.Stop()

	for {
		select {
		case <-sigCh:
			if verbose {
				fmt.Println("watch-hotplug: received stop signal; exiting")
			}
			return 0

		case <-ticker.C:
			cur := readTopologySignature()
			// If hyprctl isn't ready yet (cur == ""), just continue
			if cur == "" {
				continue
			}
			if cur != lastSig {
				// Topology changed
				if verbose {
					fmt.Println("watch-hotplug: change detected:", lastSigShort(lastSig), "→", lastSigShort(cur))
				}
				lastSig = cur
				// Debounce: reset/start timer
				if debouncing {
					if !debounceTimer.Stop() {
						select {
						case <-debounceTimer.C:
						default:
						}
					}
				}
				debouncing = true
				debounceTimer.Reset(time.Duration(debounceMS) * time.Millisecond)
			}

		case <-debounceTimer.C:
			debouncing = false
			// Enforce policy now that topology has settled briefly
			if verbose {
				fmt.Println("watch-hotplug: enforcing external-preferred policy")
			}
			_ = Enforce()
		}
	}
}

// readTopologySignature returns a stable signature for the current Hyprland monitor topology.
// It hashes the raw JSON from `hyprctl -j monitors`. Returns an empty string on error.
func readTopologySignature() string {
	out, err := exec.Command("hyprctl", "-j", "monitors").Output()
	if err != nil {
		return ""
	}
	// Normalize whitespace to reduce spurious hash changes
	norm := strings.TrimSpace(string(out))
	sum := sha256.Sum256([]byte(norm))
	return hex.EncodeToString(sum[:])
}

func lastSigShort(s string) string {
	if s == "" {
		return "none"
	}
	return s[:min(8, len(s))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
