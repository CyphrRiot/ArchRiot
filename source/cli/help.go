package cli

import "fmt"

// ShowHelp prints the CLI help message for the ArchRiot installer.
// This function centralizes the help text so main.go can remain
// delegation-only and small.
func ShowHelp() {
	fmt.Printf(`ArchRiot - The (Arch) Linux System You've Always Wanted

Usage:
  archriot              Run the main installer
  archriot --install    Run the main installer (explicit)
  archriot --upgrade    Launch the TUI upgrade flow
  archriot --tools      Launch optional tools interface
  archriot --validate   Validate packages.yaml configuration
  archriot --version    Show version information
  archriot --help       Show this help message

Options:
  -t, --tools                    Access optional advanced tools (Secure Boot, etc.)
  --apply-wallpaper-theme PATH   Apply dynamic theming based on wallpaper
  --toggle-dynamic-theming BOOL  Enable/disable dynamic theming (true/false)
  --strict-abi                   Block install if compositor/Wayland upgrades are pending
      --validate       Validate configuration without installing
  -v, --version        Display version information
  -h, --help           Display this help message

Commonly used CLI flags:
  --waybar-launch               Single-instance Waybar launcher (detached; logs to ~/.cache/archriot/runtime.log)
  --waybar-reload               Safe Waybar reload (SIGUSR2 first; restart on crash)
  --waybar-sweep                Sweep Waybar windows vs monitors; restart if duplicates appear
  --waybar-restart              Force restart Waybar (single-instance guarded)
  --waybar-pomodoro             Waybar Pomodoro JSON emitter (native replacement)
  --zed                         Focus-or-launch Zed (native > Flatpak; GL backend on Intel)
  --volume <subcmd>             Speaker/mic control: toggle|inc|dec|get|mic-toggle|mic-inc|mic-dec|mic-get
  --brightness <arg>            Backlight control: up|down|set <0-100>|get
  --stay-awake [command...]     Inhibit sleep while a command runs (systemd-inhibit)

Examples:
  archriot             # Start installation
  archriot --tools     # Open tools menu
  archriot --validate  # Check config for errors

For more information, visit: https://github.com/CyphrRiot/ArchRiot
`)
}
