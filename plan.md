# ArchRiot Development Plan — v3.7 (current)

Purpose

- Keep this plan accurate and actionable for a clean release.
- One change at a time → make → verify → stage → proceed (no commit until explicitly told “commit”).
- No duplicates; only what remains to ship v3.7.

Non‑negotiables (Build, Git, and Process)

- After any code change:
    - Run make and verify it completes successfully (green build).
    - Stage changes explicitly with exact paths:
        - git add path/to/added_or_modified_file
        - git rm path/to/removed_file
    - Confirm status with: git status -s
- Commit discipline:
    - Do not commit unless explicitly told: commit
    - Ensure the build is green at the moment of commit.
- PR hygiene:
    - Keep diffs minimal and scoped.
    - Remove legacy/obsolete assets from the index when consolidating or replacing code/assets.
- Never modify the ISO from README.md (special: version release).
- No environment variables for behavior; use flags instead (e.g., --workers).
- Consolidate variables/settings at a shared, visible level; avoid duplication.

Main.go discipline (non‑negotiable)

- main.go is an entrypoint only:
    - Parse flags; delegate to packages.
    - No feature logic in main.go.
    - Refactors must reduce lines in main.go or be rejected.
- Every CLI flag must map to a dedicated package function.
- Scope creep guard:
    - If a change exceeds ~50 lines in main.go, stop and extract a package first.
    - If a feature spans multiple concerns, split into packages and wire from main.go.
- Gate before merge:
    - Changes touching main.go must be minimal and only add/remove delegation lines.
    - New code in main.go is a defect unless it’s pure delegation or trivial glue.

Status: Completed for v3.7 (since reset)

- Telegram launcher (SUPER+G)
    - Hyprland binds simplified to PATH-resolved Telegram + 2s notification:
        - bind = $mod, G, exec, sh -lc 'notify-send -t 2000 "Telegram" "Opening Telegram..." >/dev/null 2>&1; Telegram'
    - README examples updated to match.
- Hyprland reload coalescer for dynamic colors
    - Implemented debounced reload (300ms) in session/reload.go
    - Integrated in theming flows:
        - theming.ApplyWallpaperTheme() calls session.ReloadHyprland()
        - theming.ToggleDynamicTheming() calls session.ReloadHyprland()
    - Removed temporary reload notifier from hyprland.conf.
- Keybindings Help behavior on reload
    - hyprland.conf generates help mapping on reload only:
        - exec = $HOME/.local/share/archriot/install/archriot --help-binds-generate
    - SUPER+SHIFT+H opens help via --help-binds-web.
- Workspaces and Waybar routing
    - Hyprland: default workspace keybinds limited to 1–6 (removed 0/7–10).
    - Waybar: on-click uses native "activate" (per-monitor context).
    - Waybar persistence: only 1–4 persist by default; 5–6 appear only when in use.
- GNOME Text Editor defaults
    - Custom font enforced to “Paper Mono 12” (use-system-font=false).
- Modularization (extractions and delegations)
    - windows.Switcher() for --switch-window
    - waybartools.SetupTemperature() for --setup-temperature
    - displays.Autogen() for --kanshi-autogen
    - session.Inhibit() for --stay-awake
    - session.PowerMenu() for --power-menu (added “Control Panel” option)
    - tools.UpgradeSmokeTest() for --upgrade-smoketest
    - session.SuspendGuard() for --suspend-if-undocked
    - session.MullvadStartup() for --mullvad-startup
    - session.PomodoroClick()/PomodoroDelayToggle() for pomodoro click handling
    - session.WorkspaceClick() for --waybar-workspace-click
    - session.WelcomeLaunch() for --welcome
    - session.StabilizeSession() for --stabilize-session [--inhibit]
    - upgradeguard.PreInstall(strictABI) replaces inline preInstallUpgradeGuard
    - secboot.RestoreHyprlandConfig() extracted; logger injection supported
    - secboot.RunContinuation() extracted for the Secure Boot continuation TUI
    - installer.EnsureYay() extracted yay bootstrapping
    - cli.ShowHelp() and cli.ValidateConfig() extracted
    - Removed unreachable/dead code in main.go for Waybar flags
    - main.go reduced from ~1466 lines to ~800 lines; entrypoint is delegation-only

Runtime validation (quick smoke, from repo binary install/archriot)

- Basics
    - ./install/archriot --version
    - ./install/archriot --help
- Waybar lifecycle
    - ./install/archriot --waybar-status # prints running/stopped
    - ./install/archriot --waybar-reload # debounced, safe reload
    - ./install/archriot --waybar-sweep # scan/clean duplicates
    - ./install/archriot --waybar-launch # guarded launch
- Workspaces
    - Verify Waybar shows only 1–4 persistently; 5–6 appear when they have windows.
    - Clicks route correctly per monitor (native “activate”).
- Theming & dynamic colors (single reload per action)
    - ./install/archriot --apply-wallpaper-theme PATH/TO/WALLPAPER
    - ./install/archriot --toggle-dynamic-theming true
    - Toggle back false; observe single reload each time.
- Apps / session
    - SUPER+G: verify 2s “Opening Telegram…” then Telegram launches/focuses.
    - ./install/archriot --power-menu (Control Panel present; launches)
    - ./install/archriot --stay-awake sleep 5
    - ./install/archriot --mullvad-startup (minimized; conditional)
    - ./install/archriot --stabilize-session [--inhibit]
    - ./install/archriot --welcome (if present)
- Tools
    - ./install/archriot --switch-window (hyprctl + fuzzel)
    - ./install/archriot --setup-temperature (no-op if “hwmon-path” key absent)
    - ./install/archriot --kanshi-autogen (writes ~/.config/kanshi/config)
    - ./install/archriot --upgrade-smoketest --json --quiet
- Help & docs
    - Hyprland reload triggers: archriot --help-binds-generate only (no auto-open)
    - SUPER+SHIFT+H launches --help-binds-web

Remaining scope for v3.7 (release blockers)

- Documentation & QA
    - Refresh README sections to match:
        - Telegram binding behavior (PATH + 2s notification)
        - Workspace persistence (1–4), dynamic 5–6
        - Waybar on-click routing with native “activate”
        - Updated CLI reference (extracted flags now present in respective packages)
    - Update QA matrix with the above runtime validation checklist.
- Memory defaults
    - Confirm safe, conservative memory defaults by default; advanced tuning remains opt-in (clarify docs).
- Launcher evaluation (feature-gated)
    - Evaluate Hyprlauncher as a Fuzzel replacement.
    - If parity and stability are perfect, gate behind a feature flag (default remains Fuzzel).
- Cleanup and consistency
    - Confirm removal of obsolete scripts from index (already staged):
        - config/bin/suspend-if-undocked.sh
        - config/bin/scripts/generate_keybindings_help.sh(.old), waybar-memory-accurate.py, waybar-tomato-timer.py
    - Ensure no remaining unreachable or duplicate code paths.
- Final validation pass
    - Run through the full runtime validation above on at least one real multi-monitor system and a single-monitor laptop (lid open/close and USB‑C hotplug).

Release checklist (v3.7)

- Build/verify:
    - make
    - ./install/archriot --version
- README badge/version aligned.
- Tag and push (when approved by “commit” gate):
    - git commit -am "Release v3.7: [high‑level notes]"
    - git tag -a v3.7 -m "Version 3.7 release: [details]"
    - git push origin master
    - git push origin v3.7
- Special: Do not modify the ISO reference in README.md.

Action queue (one change at a time, then make, then stage)

1. Docs & QA
    - Update README to reflect all finalized runtime behaviors (workspaces, Telegram, on-click, help flows).
    - Add a concise QA checklist (copy from “Runtime validation”) into docs.
2. Memory defaults
    - Document defaults; ensure code aligns; leave advanced tuning opt-in.
3. Hyprlauncher evaluation (feature-gated)
    - Add feature flag; wire alternative launcher; test parity; default remains Fuzzel.
4. Final test pass
    - Run the runtime validation on a real multi-monitor setup and a single-monitor laptop.
5. Prepare release
    - Follow the release checklist precisely; only commit when explicitly told: commit.

Guardrails (to not screw up again)

- Only one change in scope at a time; run make after every change.
- Stage exact files; re-check git status -s before proceeding.
- Keep main.go thin; extract logic as soon as it grows.
- No sneaky environment variables; surface flags instead.
- No orphaned code; remove dead/unreachable code immediately.
- Do not modify ISO bits in README.md.
- If you are not 100% confident a change won’t have side effects: do not propose the change yet—evaluate and test first.

Appendix: Quick commands

- Build: make
- Version: ./install/archriot --version
- Help: ./install/archriot --help
- Secure Boot continuation (TUI): handled via secboot; logger is injected.
- Workspace click (custom modules): ./install/archriot --waybar-workspace-click {name}
- Waybar reload: ./install/archriot --waybar-reload
- Tests: see “Runtime validation” above.
