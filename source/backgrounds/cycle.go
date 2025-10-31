package backgrounds

// Package backgrounds provides wallpaper management utilities.
// This file implements cycling to the next wallpaper via swaybg,
// extracted from the legacy "--swaybg-next" CLI behavior.
//
// Behavior (Next):
// - Enumerate supported images in ~/.local/share/archriot/backgrounds
// - Read the current selection from ~/.config/archriot/.current-background
// - Advance to the next image (sorted order, wraps around)
// - Persist selection and restart swaybg detached with fill mode
// - Best-effort: trigger dynamic theming via "--apply-wallpaper-theme" on the archriot binary
// - Prints a short status message; returns 0 on success, non-zero on error
//
// This function never exits the process (no os.Exit).
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

// Next cycles to the next wallpaper image and restarts swaybg.
// It returns 0 on success, non-zero otherwise.
func Next() int {
	home := os.Getenv("HOME")
	bgsDir := filepath.Join(home, ".local", "share", "archriot", "backgrounds")
	stateFile := filepath.Join(home, ".config", "archriot", ".current-background")

	entries, err := os.ReadDir(bgsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Backgrounds directory not found: %s\n", bgsDir)
		return 1
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		lower := strings.ToLower(name)
		if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") ||
			strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".webp") {
			files = append(files, filepath.Join(bgsDir, name))
		}
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No background images found in %s\n", bgsDir)
		return 1
	}
	sort.Strings(files)

	current := ""
	if b, err := os.ReadFile(stateFile); err == nil {
		current = strings.TrimSpace(string(b))
	}
	idx := -1
	for i, f := range files {
		if f == current {
			idx = i
			break
		}
	}
	next := files[(idx+1)%len(files)]

	_ = os.MkdirAll(filepath.Dir(stateFile), 0o755)
	_ = os.WriteFile(stateFile, []byte(next+"\n"), 0o644)

	// Restart swaybg
	_ = exec.Command("pkill", "-x", "swaybg").Run()
	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command("swaybg", "-i", next, "-m", "fill")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Start()

	// Best-effort theme refresh
	if self, err := os.Executable(); err == nil {
		_ = exec.Command(self, "--apply-wallpaper-theme", next).Start()
	}

	time.Sleep(1 * time.Second)
	if exec.Command("pgrep", "-x", "swaybg").Run() == nil {
		fmt.Printf("🖼️  Switched to background: %s\n", filepath.Base(next))
		fmt.Println("✓ Background service restarted")
		return 0
	}

	fmt.Println("⚠ Background service may not have started properly")
	return 0
}
