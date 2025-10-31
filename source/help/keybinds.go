package help

// Package help implements Keybindings Help generation, printing, and opening.
//
// Canonical behavior:
// - Single output path: ~/.cache/archriot/help/keybindings.html
// - Include only "bind" and "bindm" lines that use SUPER (i.e., contain $mod or SUPER)
// - Honor "unbind" lines from keybindings.conf
// - Responsive HTML: width adapts to viewport, bind column does not wrap

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PrintBinds prints hyprland keybindings (bind/bindm with SUPER/$mod only) to stdout.
// If filter is non-empty, only lines containing the filter (case-insensitive) are printed.
// Returns 0 on success, 1 on error.
func PrintBinds(filter string) int {
	home := os.Getenv("HOME")
	path := filepath.Join(home, ".config", "hypr", "hyprland.conf")
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read %s: %v\n", path, err)
		return 1
	}

	filt := strings.ToLower(strings.TrimSpace(filter))
	sc := bufio.NewScanner(strings.NewReader(string(b)))
	fmt.Println("Hyprland keybindings (bind; SUPER only):")
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)
		if !strings.HasPrefix(trim, "bind") {
			continue
		}
		cmd := firstToken(trim)
		// Exclude mouse binds (bindm/mouse:*)
		if cmd != "bind" {
			continue
		}
		if !hasSuper(trim) {
			continue
		}
		if isMouseBind(trim) {
			continue
		}
		if filt == "" || strings.Contains(strings.ToLower(line), filt) {
			fmt.Println(line)
		}
	}
	return 0
}

// GenerateHTMLAndPrintPath generates the keybindings HTML and prints the path.
// Returns 0 on success, 1 on failure.
func GenerateHTMLAndPrintPath() int {
	p, err := GenerateHTML()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	fmt.Println(p)
	return 0
}

// OpenWeb generates (if needed) and opens the keybindings HTML in a browser app window.
// Prefers brave/brave-browser (with a stable class), then falls back to xdg-open.
func OpenWeb() error {
	p, err := GenerateHTML()
	if err != nil {
		return err
	}
	url := "file://" + p

	// Prefer brave, then brave-browser; fallback to xdg-open
	class := "brave-archriot-keybinds"
	if _, err := exec.LookPath("brave"); err == nil {
		return exec.Command("brave", "--start-maximized", "--app="+url, "--class="+class).Start()
	}
	if _, err := exec.LookPath("brave-browser"); err == nil {
		return exec.Command("brave-browser", "--start-maximized", "--app="+url, "--class="+class).Start()
	}
	if _, err := exec.LookPath("xdg-open"); err == nil {
		return exec.Command("xdg-open", url).Start()
	}
	return fmt.Errorf("no suitable opener found (brave/brave-browser/xdg-open)")
}

// GenerateHTML builds the single canonical Keybindings Help page and writes it to:
// ~/.cache/archriot/help/keybindings.html
// It returns the full path or an error.
func GenerateHTML() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME not set")
	}
	hc := filepath.Join(home, ".config", "hypr", "hyprland.conf")
	uc := filepath.Join(home, ".config", "hypr", "keybindings.conf")

	outDir := filepath.Join(home, ".cache", "archriot", "help")
	outFile := filepath.Join(outDir, "keybindings.html")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create output dir: %w", err)
	}

	type row struct {
		Bind string
		Desc string
	}
	collected := map[string]row{}
	srcs := []string{}

	// Build unbind set from user config (keybindings.conf)
	unbound := map[string]struct{}{}
	if b, err := os.ReadFile(uc); err == nil {
		sc := bufio.NewScanner(strings.NewReader(string(b)))
		for sc.Scan() {
			ln := strings.TrimSpace(sc.Text())
			if !strings.HasPrefix(ln, "unbind") {
				continue
			}
			// Strip inline comment
			if idx := strings.Index(ln, "#"); idx >= 0 {
				ln = strings.TrimSpace(ln[:idx])
			}
			ln = strings.TrimSpace(strings.TrimPrefix(ln, "unbind"))
			ln = strings.TrimSpace(strings.TrimPrefix(ln, "="))
			if ln == "" {
				continue
			}
			if key, ok := normalizeKey("bind = " + ln); ok {
				unbound[key] = struct{}{}
			}
		}
	}

	// Collect from hyprland.conf first
	if b, err := os.ReadFile(hc); err == nil {
		srcs = append(srcs, hc)
		sc := bufio.NewScanner(strings.NewReader(string(b)))
		for sc.Scan() {
			ln := sc.Text()
			trim := strings.TrimSpace(ln)
			if !strings.HasPrefix(trim, "bind") {
				continue
			}
			cmd := firstToken(trim)
			// Exclude mouse binds (bindm/mouse:*)
			if cmd != "bind" {
				continue
			}
			if !hasSuper(trim) {
				continue
			}
			if isMouseBind(trim) {
				continue
			}
			desc := ""
			if idx := strings.Index(trim, "#"); idx >= 0 {
				desc = strings.TrimSpace(trim[idx+1:])
				trim = strings.TrimSpace(trim[:idx])
			}
			if key, ok := normalizeKey(trim); ok {
				if _, killed := unbound[key]; killed {
					continue
				}
				// Prefer first occurrence from base config
				if _, exists := collected[key]; !exists {
					collected[key] = row{Bind: displayBind(trim), Desc: desc}
				}
			}
		}
	}

	// Overlay from keybindings.conf (user file)
	if b, err := os.ReadFile(uc); err == nil {
		srcs = append(srcs, uc)
		sc := bufio.NewScanner(strings.NewReader(string(b)))
		for sc.Scan() {
			ln := sc.Text()
			trim := strings.TrimSpace(ln)
			if !strings.HasPrefix(trim, "bind") {
				continue
			}
			cmd := firstToken(trim)
			// Exclude mouse binds (bindm/mouse:*)
			if cmd != "bind" {
				continue
			}
			if !hasSuper(trim) {
				continue
			}
			if isMouseBind(trim) {
				continue
			}
			desc := ""
			if idx := strings.Index(trim, "#"); idx >= 0 {
				desc = strings.TrimSpace(trim[idx+1:])
				trim = strings.TrimSpace(trim[:idx])
			}
			if key, ok := normalizeKey(trim); ok {
				if _, killed := unbound[key]; killed {
					continue
				}
				collected[key] = row{Bind: displayBind(trim), Desc: desc}
			}
		}
	}

	// Flatten map to rows (order not guaranteed; stable ordering is not required)
	var rows []row
	for _, r := range collected {
		rows = append(rows, r)
	}

	// Render HTML
	var html strings.Builder
	html.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\">")
	html.WriteString("<title>ArchRiot — Keybindings Help</title>")
	html.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	// Responsive width; bind does not wrap
	html.WriteString("<style>:root{--bg:#0f172a;--fg:#cbd5e1;--muted:#94a3b8;--accent:#7aa2f7;--line:#22263a}")
	html.WriteString("html,body{background:var(--bg);color:var(--fg);margin:0;padding:0;overflow-x:hidden;font-family:ui-monospace,monospace}")
	html.WriteString(".wrap{width:96vw;margin:0 auto;padding:28px 20px 40px}")
	html.WriteString("h1{color:var(--accent);font-size:28px;margin:0 0 14px;line-height:1.15}")
	html.WriteString("table{width:100%;border-collapse:collapse;margin:12px 0}")
	html.WriteString("th,td{padding:8px 10px;border-bottom:1px solid var(--line);vertical-align:top}")
	html.WriteString("th{text-align:left}.bind{white-space:pre;overflow:visible}")
	html.WriteString("</style>")
	html.WriteString("</head><body><div class=\"wrap\">")
	html.WriteString("<h1>ArchRiot — Keybindings Help</h1>")

	html.WriteString("<table><thead><tr><th>Bind</th><th>Description</th></tr></thead><tbody>")
	esc := func(s string) string {
		return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
	}
	for _, r := range rows {
		html.WriteString("<tr><td class=\"bind\">" + esc(r.Bind) + "</td><td>" + esc(r.Desc) + "</td></tr>")
	}
	html.WriteString("</tbody></table>")

	// Sources footer
	if len(srcs) > 0 {
		html.WriteString("<p>Sources:</p><ul>")
		for _, s := range srcs {
			html.WriteString("<li>" + esc(s) + "</li>")
		}
		html.WriteString("</ul>")
	} else {
		html.WriteString("<p>No Hyprland config found. Checked:</p><ul>")
		html.WriteString("<li>" + esc(hc) + "</li>")
		html.WriteString("<li>" + esc(uc) + "</li>")
		html.WriteString("</ul>")
	}
	html.WriteString("</div></body></html>")

	if err := os.WriteFile(outFile, []byte(html.String()), 0o644); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", outFile, err)
	}
	return outFile, nil
}

// firstToken returns the first token (until space or '=') of the line in lowercase.
func firstToken(s string) string {
	s = strings.TrimSpace(s)
	i := strings.IndexAny(s, " =")
	if i < 0 {
		return strings.ToLower(s)
	}
	return strings.ToLower(strings.TrimSpace(s[:i]))
}

// hasSuper returns true if the line indicates SUPER usage via $mod or "super" (case-insensitive).
func hasSuper(line string) bool {
	l := strings.ToLower(line)
	return strings.Contains(l, "$mod") || strings.Contains(l, "super")
}

// isMouseBind returns true if the bind line’s second CSV field is a mouse binding (mouse:*).
func isMouseBind(line string) bool {
	s := strings.TrimSpace(line)
	// strip command token and optional "="
	if i := strings.Index(s, " "); i >= 0 {
		s = strings.TrimSpace(s[i+1:])
	}
	if strings.HasPrefix(s, "=") {
		s = strings.TrimSpace(s[1:])
	}
	parts := splitCSV(s)
	if len(parts) < 2 {
		return false
	}
	return strings.Contains(strings.ToLower(strings.TrimSpace(parts[1])), "mouse:")
}

// normalizeKey parses a "bind/bindm ..." line and returns a canonical "MOD, KEY" identifier.
// It extracts the fields after "bind" (and optional "="), splits by commas, and uses the first
// two fields (modifiers and key). It normalizes $mod → SUPER and uppercases for stable matching.
func normalizeKey(line string) (string, bool) {
	// Strip command token and optional "=", keep the remainder
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(strings.ToLower(s), "bind") {
		return "", false
	}
	// Remove leading token and equal sign if present
	rest := s
	if i := strings.Index(rest, " "); i >= 0 {
		rest = strings.TrimSpace(rest[i+1:])
	}
	if strings.HasPrefix(rest, "=") {
		rest = strings.TrimSpace(rest[1:])
	}

	// Now rest should be CSV-like: "$mod, H, exec, ..."
	parts := splitCSV(rest)
	if len(parts) < 2 {
		return "", false
	}
	mod := strings.TrimSpace(parts[0])
	key := strings.TrimSpace(parts[1])

	mod = strings.ReplaceAll(mod, "$mod", "SUPER")
	mod = strings.ReplaceAll(strings.ToUpper(mod), "  ", " ")
	key = strings.ToUpper(key)

	return mod + ", " + key, true
}

// displayBind creates a human-friendly bind display from a raw "bind = ..." line.
// It converts "$mod, H, ..." into "SUPER + H, ..." (only first two components are guaranteed).
func displayBind(line string) string {
	// Preserve the original structure but normalize $mod → SUPER and join two first components.
	s := strings.TrimSpace(line)
	rest := s
	if i := strings.Index(rest, " "); i >= 0 {
		rest = strings.TrimSpace(rest[i+1:])
	}
	if strings.HasPrefix(rest, "=") {
		rest = strings.TrimSpace(rest[1:])
	}
	parts := splitCSV(rest)
	if len(parts) == 0 {
		return s
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if len(parts) >= 2 {
		m := strings.ReplaceAll(parts[0], "$mod", "SUPER")
		return strings.TrimSpace(strings.ReplaceAll(m, ",", "")) + " + " + parts[1]
	}
	return strings.ReplaceAll(parts[0], "$mod", "SUPER")
}

// splitCSV splits on commas without special quoting rules (sufficient for hypr syntax).
func splitCSV(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}
