package tools

// UpgradeSmokeTest implements the --upgrade-smoketest logic as a reusable function.
// It compares the desired packages (from packages.yaml) against currently installed
// packages (via pacman -Qq) to detect potential reintroductions. It supports:
//
// Flags:
//   --config PATH   Path to packages.yaml (default: $HOME/.local/share/archriot/install/packages.yaml)
//   --json          Emit JSON output
//   --quiet         Suppress human-friendly messages (non-JSON)
//   -h, --help      Show usage
//
// Exit codes:
//   0 -> OK (no potential reintroductions)
//   2 -> Potential reintroductions detected
//   3 -> Unavailable/missing prerequisites (e.g., pacman/config missing)
//
// This function does not os.Exit; it returns the intended exit code.
import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"archriot-installer/config"
)

// UpgradeSmokeTest runs the local upgrade smoke test using provided args.
func UpgradeSmokeTest(args []string) int {
	home := os.Getenv("HOME")
	configPath := filepath.Join(home, ".local", "share", "archriot", "install", "packages.yaml")
	allowlistPath := filepath.Join(home, ".config", "archriot", "upgrade-allowlist.txt")
	outputJSON := false
	quiet := false

	// Parse flags
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--config":
			if i+1 < len(args) {
				configPath = args[i+1]
				i++
			} else {
				if outputJSON {
					payload := map[string]interface{}{
						"status":    "unavailable",
						"message":   "Missing value for --config",
						"config":    configPath,
						"allowlist": allowlistPath,
						"missing":   []string{},
					}
					if b, e := json.Marshal(payload); e == nil {
						fmt.Println(string(b))
					} else {
						fmt.Println(`{"status":"unavailable","message":"Missing value for --config","missing":[]}`)
					}
				} else if !quiet {
					fmt.Fprintln(os.Stderr, "Missing value for --config")
				}
				return 3
			}
		case "--json":
			outputJSON = true
		case "--quiet":
			quiet = true
		case "-h", "--help":
			if !outputJSON && !quiet {
				fmt.Println("ArchRiot Local Upgrade Smoke Test")
				fmt.Println("Usage: archriot --upgrade-smoketest [--config PATH] [--json] [--quiet]")
			}
			return 0
		default:
			if !quiet {
				fmt.Fprintf(os.Stderr, "Unknown argument: %s\n", arg)
			}
			return 3
		}
	}

	// Ensure pacman exists
	if _, err := exec.LookPath("pacman"); err != nil {
		if outputJSON {
			payload := map[string]interface{}{
				"status":    "unavailable",
				"message":   "pacman not found in PATH; cannot assess installed packages",
				"config":    configPath,
				"allowlist": allowlistPath,
				"missing":   []string{},
			}
			if b, e := json.Marshal(payload); e == nil {
				fmt.Println(string(b))
			} else {
				fmt.Println(`{"status":"unavailable","message":"pacman not found in PATH","missing":[]}`)
			}
		} else if !quiet {
			fmt.Fprintln(os.Stderr, "pacman not found; cannot assess installed packages.")
		}
		return 3
	}

	// Ensure config exists
	if _, err := os.Stat(configPath); err != nil {
		if outputJSON {
			payload := map[string]interface{}{
				"status":    "unavailable",
				"message":   fmt.Sprintf("packages.yaml not found at: %s", configPath),
				"config":    configPath,
				"allowlist": allowlistPath,
				"missing":   []string{},
			}
			if b, e := json.Marshal(payload); e == nil {
				fmt.Println(string(b))
			} else {
				fmt.Printf("{\"status\":\"unavailable\",\"message\":\"packages.yaml not found at: %s\",\"missing\":[]}\n", configPath)
			}
		} else if !quiet {
			fmt.Fprintf(os.Stderr, "packages.yaml not found at: %s\n", configPath)
		}
		return 3
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		if outputJSON {
			payload := map[string]interface{}{
				"status":    "unavailable",
				"message":   fmt.Sprintf("failed to load packages.yaml: %v", err),
				"config":    configPath,
				"allowlist": allowlistPath,
				"missing":   []string{},
			}
			if b, e := json.Marshal(payload); e == nil {
				fmt.Println(string(b))
			} else {
				fmt.Printf("{\"status\":\"unavailable\",\"message\":\"failed to load packages.yaml\",\"missing\":[]}\n")
			}
		} else if !quiet {
			fmt.Fprintf(os.Stderr, "failed to load packages.yaml: %v\n", err)
		}
		return 3
	}

	// Collect desired packages from all module maps via reflection
	desired := make(map[string]struct{})
	cv := reflect.ValueOf(cfg)
	if cv.Kind() == reflect.Ptr {
		cv = cv.Elem()
	}
	ct := cv.Type()
	for i := 0; i < cv.NumField(); i++ {
		field := cv.Field(i)
		// Only consider map[string]config.Module fields
		if field.Kind() != reflect.Map || field.Type().String() != "map[string]config.Module" {
			continue
		}
		iter := field.MapRange()
		for iter.Next() {
			modVal := iter.Value()
			mod, ok := modVal.Interface().(config.Module)
			if !ok {
				continue
			}
			for _, p := range mod.Packages {
				if p = strings.TrimSpace(p); p != "" {
					desired[p] = struct{}{}
				}
			}
		}
	}
	_ = ct // keep ct in case of future use or build tags

	// Load allowlist (optional)
	allow := make(map[string]struct{})
	if f, err := os.Open(allowlistPath); err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := sc.Text()
			if idx := strings.IndexByte(line, '#'); idx >= 0 {
				line = line[:idx]
			}
			line = strings.TrimSpace(line)
			if line != "" {
				allow[line] = struct{}{}
			}
		}
		_ = f.Close()
	}

	// Get installed packages
	out, err := exec.Command("pacman", "-Qq").Output()
	if err != nil {
		if outputJSON {
			payload := map[string]interface{}{
				"status":    "unavailable",
				"message":   "failed to query installed packages with pacman -Qq",
				"config":    configPath,
				"allowlist": allowlistPath,
				"missing":   []string{},
			}
			if b, e := json.Marshal(payload); e == nil {
				fmt.Println(string(b))
			} else {
				fmt.Println(`{"status":"unavailable","message":"failed to query installed packages","missing":[]}`)
			}
		} else if !quiet {
			fmt.Fprintln(os.Stderr, "failed to query installed packages with pacman -Qq")
		}
		return 3
	}
	installed := make(map[string]struct{})
	{
		sc := bufio.NewScanner(strings.NewReader(string(out)))
		for sc.Scan() {
			p := strings.TrimSpace(sc.Text())
			if p != "" {
				installed[p] = struct{}{}
			}
		}
	}

	// Compute missing (present in YAML but not installed, and not allowlisted)
	var missing []string
	for p := range desired {
		if _, ok := allow[p]; ok {
			continue
		}
		if _, ok := installed[p]; !ok {
			missing = append(missing, p)
		}
	}

	// Output
	if outputJSON {
		status := "ok"
		message := "No potential reintroductions detected"
		if len(missing) > 0 {
			status = "warn"
			message = "Potential reintroductions detected"
		}
		payload := map[string]interface{}{
			"status":    status,
			"message":   message,
			"config":    configPath,
			"allowlist": allowlistPath,
			"missing":   missing,
		}
		if b, e := json.Marshal(payload); e == nil {
			fmt.Println(string(b))
		} else {
			// Minimal fallback
			fmt.Printf("{\"status\":\"%s\",\"message\":\"%s\",\"missing_count\":%d}\n", status, message, len(missing))
		}
	} else if !quiet {
		fmt.Println("ArchRiot Local Upgrade Smoke Test")
		fmt.Printf("Config:    %s\n", configPath)
		if _, err := os.Stat(allowlistPath); err == nil {
			fmt.Printf("Allowlist: %s\n", allowlistPath)
		}
		fmt.Println()
		if len(missing) == 0 {
			fmt.Println("✅ No potential reintroductions detected.")
			fmt.Println("Local upgrade appears safe with respect to previously removed packages.")
		} else {
			fmt.Println("⚠️  Potential reintroductions detected (present in packages.yaml but not installed):")
			for _, p := range missing {
				fmt.Printf("  - %s\n", p)
			}
			fmt.Println()
			fmt.Println("This suggests these packages were removed (or never installed) and would be installed by a Local upgrade.")
			if _, err := os.Stat(allowlistPath); err == nil {
				fmt.Printf("You can add specific packages to %s (one per line) to suppress this warning.\n", allowlistPath)
			}
		}
	}

	// Exit codes
	if len(missing) == 0 {
		return 0
	}
	return 2
}
