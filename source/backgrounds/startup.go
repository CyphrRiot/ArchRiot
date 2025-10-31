package backgrounds

// Package backgrounds provides wallpaper management utilities.
// This file implements starting the wallpaper service at login
// using a previously selected background (with sensible fallbacks),
// extracted from the legacy "--startup-background" behavior.
//
// Behavior (Startup):
// - Read ~/.config/archriot/background-prefs.json and extract "current_background"
// - Fallback to a default name ("riot_01.jpg") or the first image in the backgrounds dir
// - Persist the selection to ~/.config/archriot/.current-background
// - Restart swaybg detached with fill mode
// - Do NOT apply dynamic theming at startup (to avoid races)
// - Returns 0 on success (best-effort), non-zero on fatal errors
//
// No os.Exit calls; the caller handles process exit codes.
import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Startup picks a wallpaper at login based on saved preferences and starts swaybg.
// It avoids applying dynamic theming at startup to reduce race conditions.
// Returns 0 on success; non-zero on error.
func Startup() int {
	home := os.Getenv("HOME")
	configFile := filepath.Join(home, ".config", "archriot", "background-prefs.json")
	bgsDir := filepath.Join(home, ".local", "share", "archriot", "backgrounds")
	stateFile := filepath.Join(home, ".config", "archriot", ".current-background")
	defaultName := "riot_01.jpg"

	// Read desired background name from JSON (current_background)
	bgName := ""
	if b, err := os.ReadFile(configFile); err == nil {
		var obj map[string]interface{}
		if json.Unmarshal(b, &obj) == nil {
			if v, ok := obj["current_background"]; ok {
				if s, ok2 := v.(string); ok2 {
					bgName = s
				}
			}
		}
	}
	if strings.TrimSpace(bgName) == "" || strings.EqualFold(bgName, "null") {
		bgName = defaultName
	}

	// Resolve file path with fallbacks
	pick := filepath.Join(bgsDir, bgName)
	if st, err := os.Stat(pick); err != nil || st.IsDir() {
		pick = filepath.Join(bgsDir, defaultName)
	}
	if st, err := os.Stat(pick); err != nil || st.IsDir() {
		// Fallback: pick the first supported image in the dir
		entries, _ := os.ReadDir(bgsDir)
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			lower := strings.ToLower(name)
			if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") ||
				strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".webp") {
				pick = filepath.Join(bgsDir, name)
				break
			}
		}
	}

	if pick == "" {
		// No backgrounds found; nothing to do
		return 0
	}

	// Update state file for runtime cycling compatibility
	_ = os.MkdirAll(filepath.Dir(stateFile), 0o755)
	_ = os.WriteFile(stateFile, []byte(pick+"\n"), 0o644)

	// Relaunch swaybg (detached); do not apply theme on startup to avoid conflicts
	_ = exec.Command("pkill", "-x", "swaybg").Run()
	time.Sleep(300 * time.Millisecond)

	cmd := exec.Command("swaybg", "-i", pick, "-m", "fill")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	_ = cmd.Start()

	time.Sleep(500 * time.Millisecond)
	return 0
}
