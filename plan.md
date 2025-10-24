# ArchRiot Development Plan — Refreshed (v3.6.1)

Purpose

- Keep this plan accurate and actionable.
- Remove completed items from open tasks.
- Capture newly completed work and define focused next steps.
- Maintain a strict one-change-at-a-time cadence: propose → edit → make → verify.

---

Completed in v3.6.1 (shipped)

- Zed stability: `archriot --zed` prefers `WGPU_BACKEND=gl` on Intel; Hyprland bind updated to use the hardened launcher.
- Telegram launcher: native-first, broadened class matching, and runtime logging for focus/launch reliability.
- Control Panel display scaling: focused-monitor parsing; robust `hyprctl` apply; kanshi config best-effort update.
- Mullvad toggle semantics: manage `exec-once` `archriot --mullvad-startup` and respect auto-connect; `--mullvad-startup` skips GUI when auto-connect is off.
- Waybar UX: tooltip opacity increased for calendar readability; native Pomodoro JSON emitter; module wired to `archriot --waybar-pomodoro`.
- Window recovery: improved off-screen window fixer with multi-monitor awareness and notifications.
- Bootstrap cleanup: remove ISO one-shot from `~/.bashrc` and stop adding it during install.
- Docs/version: `VERSION` and README badge/current release updated to v3.6.1.

Completed in v3.6 (shipped)

- Installer robustness: per‑command timeouts with output capture to avoid hangs; non‑critical commands continue, critical ones fail fast; extended timeouts for pacman/yay.
- Control Panel (GTK) safety: reapply only when a graphical session is detected; force GSK_RENDERER=gl; automatic cairo software‑renderer fallback if GL fails.
- Hypridle reliability: use absolute path in autostart (`/usr/bin/hypridle`) and ensureIdleLockSetup appends the same to avoid PATH/env issues post‑upgrade.
- Installer hardening: make `hypr-dock-inhibit.service` disable idempotent and non‑fatal (`|| true`) when the unit is missing.
- README/docs polish: fixed TOC anchors and Quick Install deep links; added optional Brave multi‑monitor flags; kept “What’s New” at the end; version badge updated.
- README/docs polish (more): removed duplicate “Memory Tuning (Opt‑in)” section; appended “What’s New in v3.6” at bottom.
- Help windows: widened SUPER+H and SUPER+SHIFT+H help windows by 15% to reduce wrapping.
- Version bump: v3.6 shipped.

Completed in v3.5.1 (shipped)

- README.md overhaul: Navigation/“📚 Navigate This Guide”, Quick Install panel, Advanced Usage/Troubleshooting consolidation, Workspace Styles doc, Contributors section, “What’s New” moved to end
- Keybindings Help PWA: SUPER+SHIFT+H opens local help (Brave app) with stable/fallback window rules; GTK fallback available; binding updated
- Mullvad/VPN stability: Removed NetworkManager reload during install to avoid VPN drop during installation
- Preflight fixes: robust portal detection (PATH, /usr/lib, or running process); cleaned Wi‑Fi power-save diagnostics output
- Installer guards: pre-install ABI advisory and --strict-abi enforcement for compositor/Wayland packages

Next Steps (upgrade-safe; one change at a time)

1. Issue verification and closure (validate, then close)

- #31 Brave crash on workspace/screen switch: retest across GPUs; document safe-mode flag if needed
- #30 OBS “Screen Capture” present: portals functional on clean install
- #24 Thunar default terminal: verify Ghostty “Open Terminal Here” on a clean install

2. README.md finalization (post‑overhaul polish)

- Final pass to remove any duplication and normalize heading styles
- Ensure all code snippets are fenced; fix spacing/indentation regressions
- Keep “What’s New” at end; ensure “Navigate This Guide” is the only TOC

3. QA matrix and docs refresh

- Matrix: Intel/AMD/NVIDIA; single/multi-monitor; fractional scaling; en-US + one non-Latin
- Record pass/fail and remediation steps; update README with confirmed guidance/screenshots

4. HyprToolkit migration (Control Panel first; feature-flagged)

5. GTK Keybindings Help (native app; PWA shipped)

- Current: PWA (Brave app) help is shipped and bound to SUPER+SHIFT+H; GTK fallback exists
- Goal (follow‑up): Build native GTK4 help (Go) with live parsing and search; dual‑ship behind a feature flag until parity
- Scope:
    - Add a small GTK4 app (Go) or reuse Control Panel framework; scrollable/searchable; grouped by category
    - Keep CLI/PWA fallback; wire a feature flag to switch default when GTK parity is achieved
- Acceptance:
    - SUPER+SHIFT+H opens the GTK Help with live keybindings and working search when feature flag is enabled
    - Categories include: Getting Started, Window Management, Applications, Communication, Screenshots, System, Web Apps, Media
    - Theming matches Control Panel; sensible window sizing; no dependency on a browser or local HTML

- Ship dual‑path for one minor release (feature‑flagged); flip default at parity
- Maintain complete CLI/PWA fallbacks

6. Script consolidation and CLI migration (no new scripts)

- Goal: Reduce floating shell scripts and migrate behavior into first-class archriot CLI flags.
- Scope (prioritized):
    - In progress: Finalize migration of Waybar helpers to native CLI where feasible; keep JSON protocol compatibility.
    - EVALUATE: Replace `suspend-if-undocked.sh` (hypridle hook) with `archriot --suspend-if-undocked` vs keep as file
    - Remove or archive unused scripts under `config/bin/scripts/` after parity is verified.

- Acceptance:
    - No new standalone scripts introduced.
    - Script count in `config/bin/scripts` reduced by at least 50% with feature parity.
    - All migrated features are documented under “ArchRiot CLI Flags.”

7. Secure Boot (sbctl) wizard — opt-in, gated, and documented (no auto-run)

- Wizard: detect → keygen → sign → pacman hooks → pre-reboot checklist
- Post-reboot continuation validates and restores normal behavior
- Document clear rollback; verify on supported hardware
