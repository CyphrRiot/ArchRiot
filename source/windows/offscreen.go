package windows

// Package windows provides window management utilities (Hyprland).
// This file extracts the legacy "--fix-offscreen-windows" behavior into
// a reusable function that recenters floating windows that have ended up
// fully outside any monitor's visible area.
//
// Behavior (FixOffscreen):
// - Requires "hyprctl" in PATH (prints an error and returns non-zero if absent)
// - Reads monitors (JSON) and builds a combined visible region
// - Reads clients (JSON) and detects floating windows whose rect doesn't
//   intersect any monitor rect (with a small tolerance)
// - For each off-screen floating window:
//   * Switch to the window's workspace
//   * Focus the window by address
//   * Center the window
//   * Switch back to the original workspace
// - Prints one line per fixed window and a summary line at the end
// - Sends a best-effort desktop notification if "notify-send" is available
//
// The function never calls os.Exit; it returns an int exit code (0 = success).
import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// FixOffscreen recenters off-screen floating windows using Hyprland's hyprctl JSON.
func FixOffscreen() int {
	// Ensure hyprctl is available
	if _, err := exec.LookPath("hyprctl"); err != nil {
		fmt.Println("hyprctl not found in PATH")
		return 1
	}

	// Read monitors (gather absolute geometry for all monitors)
	type mon struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	monOut, err := exec.Command("hyprctl", "monitors", "-j").Output()
	if err != nil {
		fmt.Printf("Failed to read monitors: %v\n", err)
		return 1
	}
	var mons []mon
	if err := json.Unmarshal(monOut, &mons); err != nil || len(mons) == 0 {
		fmt.Println("Failed to parse monitors JSON")
		return 1
	}

	// Build screen rectangles
	type rect struct{ x0, y0, x1, y1 int }
	var screens []rect
	for _, m := range mons {
		screens = append(screens, rect{m.X, m.Y, m.X + m.Width, m.Y + m.Height})
	}

	// Read current workspace from activewindow
	type active struct {
		Workspace struct {
			ID int `json:"id"`
		} `json:"workspace"`
	}
	awOut, _ := exec.Command("hyprctl", "activewindow", "-j").Output()
	curWS := 1
	if len(awOut) > 0 {
		var aw active
		if json.Unmarshal(awOut, &aw) == nil && aw.Workspace.ID != 0 {
			curWS = aw.Workspace.ID
		}
	}

	// Read clients
	type client struct {
		Address   string `json:"address"`
		At        []int  `json:"at"`
		Size      []int  `json:"size"`
		Floating  bool   `json:"floating"`
		Workspace struct {
			ID int `json:"id"`
		} `json:"workspace"`
		Class string `json:"class"`
		Title string `json:"title"`
	}
	clOut, err := exec.Command("hyprctl", "clients", "-j").Output()
	if err != nil {
		fmt.Printf("Failed to read clients: %v\n", err)
		return 1
	}
	var cls []client
	if err := json.Unmarshal(clOut, &cls); err != nil {
		fmt.Println("Failed to parse clients JSON")
		return 1
	}

	overlaps := func(a rect, bx0, by0, bx1, by1 int) bool {
		return a.x0 < bx1 && a.x1 > bx0 && a.y0 < by1 && a.y1 > by0
	}
	isOffscreen := func(x, y, w, h int) bool {
		// add a small tolerance to catch nearly invisible positions
		tol := 8
		wx0, wy0 := x, y
		wx1, wy1 := x+w, y+h
		wx0 -= tol
		wy0 -= tol
		wx1 += tol
		wy1 += tol
		for _, s := range screens {
			if overlaps(s, wx0, wy0, wx1, wy1) {
				return false
			}
		}
		return true
	}

	fixed := 0
	for _, c := range cls {
		// Only floating windows with geometry
		if !c.Floating || len(c.At) < 2 || len(c.Size) < 2 {
			continue
		}
		x, y := c.At[0], c.At[1]
		w, h := c.Size[0], c.Size[1]

		if !isOffscreen(x, y, w, h) {
			continue
		}
		ws := c.Workspace.ID
		if ws == 0 {
			ws = curWS
		}

		// Switch to window's workspace, focus, center, return
		_ = exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", ws)).Run()
		time.Sleep(100 * time.Millisecond)
		_ = exec.Command("hyprctl", "dispatch", "focuswindow", "address:"+c.Address).Run()
		time.Sleep(120 * time.Millisecond)
		_ = exec.Command("hyprctl", "dispatch", "centerwindow").Run()
		time.Sleep(120 * time.Millisecond)
		_ = exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", curWS)).Run()

		fmt.Printf("Centered off-screen window: %s (%s) from %d,%d size %dx%d on ws %d\n", c.Class, c.Title, x, y, w, h, ws)
		fixed++
	}

	// Desktop notifications (best-effort)
	if _, err := exec.LookPath("notify-send"); err == nil {
		if fixed > 0 {
			_ = exec.Command("notify-send", "-t", "2500", "ArchRiot", fmt.Sprintf("Centered %d off-screen floating window(s)", fixed)).Start()
		} else {
			_ = exec.Command("notify-send", "-t", "2000", "ArchRiot", "No off-screen floating windows").Start()
		}
	}

	if fixed > 0 {
		fmt.Printf("Fixed %d off-screen floating windows\n", fixed)
	} else {
		fmt.Println("No off-screen floating windows found")
	}
	return 0
}
