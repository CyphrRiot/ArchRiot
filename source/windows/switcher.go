package windows

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Switcher provides the Hyprland window switcher using hyprctl JSON + fuzzel.
// Returns 0 on best‑effort completion; 1 on hard precondition failures.
func Switcher() int {
	have := func(name string) bool { _, err := exec.LookPath(name); return err == nil }
	if !have("hyprctl") {
		fmt.Fprintln(os.Stderr, "hyprctl not found in PATH")
		return 1
	}
	if !have("fuzzel") {
		fmt.Fprintln(os.Stderr, "fuzzel not found in PATH")
		return 1
	}

	type client struct {
		Address   string `json:"address"`
		Workspace struct {
			ID int `json:"id"`
		} `json:"workspace"`
		Class string `json:"class"`
		Title string `json:"title"`
	}
	clOut, err := exec.Command("hyprctl", "clients", "-j").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read clients: %v\n", err)
		return 1
	}
	var cls []client
	if err := json.Unmarshal(clOut, &cls); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse clients JSON")
		return 1
	}

	type option struct {
		Display string
		Address string
		WS      int
	}
	var opts []option
	for _, c := range cls {
		display := fmt.Sprintf("[%d] %s — %s", c.Workspace.ID, c.Title, c.Class)
		opts = append(opts, option{Display: display, Address: c.Address, WS: c.Workspace.ID})
	}
	if len(opts) == 0 {
		fmt.Println("No windows found")
		return 0
	}

	var b strings.Builder
	for i, o := range opts {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(o.Display)
	}

	cmd := exec.Command("fuzzel", "--dmenu", "--prompt=Switch to: ", "--width=60", "--lines=10")
	cmd.Stdin = strings.NewReader(b.String())
	sel, err := cmd.Output()
	if err != nil {
		return 0
	}
	choice := strings.TrimSpace(string(sel))
	if choice == "" {
		return 0
	}

	var chosen *option
	for i := range opts {
		if opts[i].Display == choice {
			chosen = &opts[i]
			break
		}
	}
	if chosen == nil {
		return 0
	}

	_ = exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", chosen.WS)).Run()
	_ = exec.Command("hyprctl", "dispatch", "focuswindow", "address:"+chosen.Address).Run()
	_ = exec.Command("hyprctl", "dispatch", "bringactivetotop").Run()
	return 0
}
