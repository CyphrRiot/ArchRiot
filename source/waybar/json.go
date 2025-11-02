package waybar

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

// EmitTemp will emit temperature status JSON (moved from --waybar-temp in main.go).
