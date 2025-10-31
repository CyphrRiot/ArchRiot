package session

// Hyprland reload coalescer (debounced single reload).
//
// This utility collapses multiple reload requests into a single
// "hyprctl reload" invocation when they occur within a short window.
// Use ReloadHyprland() for debounced/background reloads.
// Use ReloadHyprlandImmediate() to force an immediate (synchronous) reload.
//
// Typical usage:
//   // Debounced reload (fire-and-forget)
//   session.ReloadHyprland()
//
//   // Immediate reload (waits for completion)
//   if err := session.ReloadHyprlandImmediate(); err != nil { /* handle */ }
//
// By default, the debounce window is 300ms. You can adjust it at runtime with
// SetReloadDebounce(d).
//
// Notes:
// - Debounced reloads run in the background and do not return an error.
// - If a reload is already pending, additional debounced requests are merged.
// - Immediate reloads bypass the coalescer and run "hyprctl reload" directly.

import (
	"os/exec"
	"sync"
	"time"
)

var (
	reloadOnce     sync.Once
	reloadCh       chan struct{}
	reloadMu       sync.Mutex
	reloadDebounce = 300 * time.Millisecond
)

// SetReloadDebounce changes the debounce interval for debounced reloads.
// Calls with d <= 0 are ignored and keep the current debounce duration.
func SetReloadDebounce(d time.Duration) {
	if d <= 0 {
		return
	}
	reloadMu.Lock()
	reloadDebounce = d
	reloadMu.Unlock()
}

// ReloadHyprland enqueues a debounced Hyprland reload. Multiple calls within
// the debounce window are coalesced into a single "hyprctl reload" execution.
// This function is non-blocking and returns immediately.
func ReloadHyprland() {
	initReloadWorker()
	// Non-blocking send; if a request is already queued, we don't need another.
	select {
	case reloadCh <- struct{}{}:
	default:
		// A reload is already pending within the debounce window.
	}
}

// ReloadHyprlandImmediate performs a synchronous "hyprctl reload" and returns
// the command error (if any). This bypasses the coalescer.
func ReloadHyprlandImmediate() error {
	cmd := exec.Command("hyprctl", "reload")
	return cmd.Run()
}

func initReloadWorker() {
	reloadOnce.Do(func() {
		reloadCh = make(chan struct{}, 1)
		go reloadWorker()
	})
}

func reloadWorker() {
	var timer *time.Timer
	for {
		// Wait for at least one request
		<-reloadCh

		// Ensure we have a timer started/reset to the current debounce
		reloadMu.Lock()
		d := reloadDebounce
		reloadMu.Unlock()

		if timer == nil {
			timer = time.NewTimer(d)
		} else {
			if !timer.Stop() {
				// Drain timer channel if needed
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(d)
		}

		// Collect additional quick bursts and reset timer each time
	collectLoop:
		for {
			select {
			case <-reloadCh:
				// Another request came in within the window; reset timer
				reloadMu.Lock()
				d = reloadDebounce
				reloadMu.Unlock()

				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(d)

			case <-timer.C:
				// Debounce window elapsed: perform a single reload
				_ = exec.Command("hyprctl", "reload").Run()
				timer = nil
				break collectLoop
			}
		}
	}
}
