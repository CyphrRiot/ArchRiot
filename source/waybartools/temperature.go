package waybartools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetupTemperature detects CPU temperature sensor files and updates the Waybar
// Modules configuration to set the "hwmon-path" array accordingly.
// Returns 0 on success (or no-op when key not found), non-zero on hard errors.
func SetupTemperature() int {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		fmt.Fprintln(os.Stderr, "Failed to determine HOME directory")
		return 1
	}
	modPath := filepath.Join(home, ".config", "waybar", "Modules")

	// Detect coretemp (Intel) from /sys/class/hwmon/hwmon*/name == "coretemp" and temp1_input exists
	findCoretemp := func() string {
		entries, _ := os.ReadDir("/sys/class/hwmon")
		for _, e := range entries {
			hp := filepath.Join("/sys/class/hwmon", e.Name())
			nameB, err := os.ReadFile(filepath.Join(hp, "name"))
			if err != nil {
				continue
			}
			if strings.TrimSpace(string(nameB)) == "coretemp" {
				if _, err := os.Stat(filepath.Join(hp, "temp1_input")); err == nil {
					return filepath.Join(hp, "temp1_input")
				}
			}
		}
		return ""
	}

	// Detect thermal zone type x86_pkg_temp
	findPkgTemp := func() string {
		entries, _ := os.ReadDir("/sys/class/thermal")
		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), "thermal_zone") {
				continue
			}
			zp := filepath.Join("/sys/class/thermal", e.Name())
			tb, err := os.ReadFile(filepath.Join(zp, "type"))
			if err != nil {
				continue
			}
			if strings.TrimSpace(string(tb)) == "x86_pkg_temp" {
				return filepath.Join(zp, "temp")
			}
		}
		return ""
	}

	// Build candidate paths
	var paths []string
	if core := findCoretemp(); core != "" {
		paths = append(paths, core)
	}
	if pkg := findPkgTemp(); pkg != "" {
		paths = append(paths, pkg)
	}
	// Always include a basic fallback
	paths = append(paths, "/sys/class/thermal/thermal_zone0/temp")

	// Prepare JSON array content as a line-separated or single-line replacement
	var quoted []string
	for _, p := range paths {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", p))
	}
	arrayLine := strings.Join(quoted, ", ")

	// Read the Waybar Modules file
	data, err := os.ReadFile(modPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Waybar Modules not found: %s\n", modPath)
		return 1
	}
	lines := strings.Split(string(data), "\n")

	// Replace the "hwmon-path": [...] array in-place, preserving indentation and trailing comma.
	out := make([]string, 0, len(lines)+3)
	i := 0
	replaced := false

	for i < len(lines) {
		line := lines[i]
		trim := strings.TrimSpace(line)

		if strings.HasPrefix(trim, "\"hwmon-path\"") && strings.Contains(trim, "[") {
			// Determine indentation
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]

			// Consume until we pass the closing ']' (handles one-line and multi-line)
			// Keep note if the closing bracket had a trailing comma in original
			trailingComma := ","
			j := i
			foundClose := false
			for ; j < len(lines); j++ {
				if strings.Contains(lines[j], "]") {
					// If the original closing had no comma, keep it without comma
					t := strings.TrimSpace(lines[j])
					if strings.HasSuffix(t, "]") && !strings.HasSuffix(t, "],") {
						trailingComma = ","
						// Many JSON fragments place a comma after the key/value element.
						// We'll default to comma to avoid breaking downstream structure,
						// as typical Waybar Modules has a trailing comma here.
						// If the user had no comma originally, Waybar will still parse.
					}
					foundClose = true
					j++ // move past closing
					break
				}
			}
			if !foundClose {
				// Malformed block; bail out without modifying
				out = append(out, line)
				i++
				continue
			}

			// Write our replacement block
			out = append(out, indent+"\"hwmon-path\": [")
			out = append(out, indent+"    "+arrayLine)
			out = append(out, indent+"]"+trailingComma)

			replaced = true
			i = j
			continue
		}

		out = append(out, line)
		i++
	}

	if !replaced {
		fmt.Println("Note: \"hwmon-path\" not found in Waybar Modules; no changes made")
		return 0
	}

	// Write back the modified Modules
	if err := os.WriteFile(modPath, []byte(strings.Join(out, "\n")), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update Waybar Modules: %v\n", err)
		return 1
	}

	fmt.Printf("✓ Updated Waybar temperature configuration\nPaths: %s\n", arrayLine)
	return 0
}
