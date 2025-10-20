# ArchRiot Development Plan — Refreshed (v3.5)

Purpose

- Keep this plan accurate and actionable.
- Remove completed items from open tasks.
- Capture newly completed work and define focused next steps.
- Maintain a strict one-change-at-a-time cadence: propose → edit → make → verify.

---

Next Steps (upgrade-safe; one change at a time)

2. Issue verification and closure (validate, then close)

- #31 Brave crash on workspace/screen switch: retest across GPUs; document safe-mode flag if needed
- #30 OBS “Screen Capture” present: portals functional on clean install
- #24 Thunar default terminal: verify Ghostty “Open Terminal Here” on a clean install

4. QA matrix and docs refresh

- Matrix: Intel/AMD/NVIDIA; single/multi-monitor; fractional scaling; en-US + one non-Latin
- Record pass/fail and remediation steps; update README with confirmed guidance/screenshots

5. HyprToolkit migration (Control Panel first; feature-flagged)

6. GTK Keybindings Help (replace SUPER+SHIFT+H browser flow)

- Goal: Launch a native GTK window that renders current keybindings (parsed from hyprland.conf or a generated JSON).
- Scope:
    - Add a small GTK4 app (Go) or reuse the Control Panel framework to render a scrollable, searchable list grouped by category.
    - Replace SUPER+SHIFT+H binding to launch this GTK window directly; keep a CLI fallback (archriot --help-binds).
- Acceptance:
    - SUPER+SHIFT+H opens the GTK Help with live keybindings and working search.
    - Categories include: Getting Started, Window Management, Applications, Communication, Screenshots, System, Web Apps, Media.
    - Theming matches Control Panel; sensible window sizing; no dependency on a browser or local HTML.

- Feature-flag for one minor release (dual-ship); flip default when parity achieved
- Maintain complete CLI fallbacks

7. Script consolidation and CLI migration (no new scripts)

- Goal: Reduce floating shell scripts and migrate behavior into first-class archriot CLI flags.
- Scope (prioritized):
    - TODO: Replace `welcome` with `archriot --welcome`
    - TODO: Replace `zed-wayland` with `archriot --zed` (focus-or-launch Zed, Wayland class handling)
    - TODO: Replace `waybar-*.py` helpers with native CLI metrics where feasible; keep JSON protocol compatibility.
    - EVALUATE: Replace `suspend-if-undocked.sh` (hypridle hook) with `archriot --suspend-if-undocked` vs keep as file
    - Remove or archive unused scripts under `config/bin/scripts/` after parity is verified.

- Acceptance:
    - No new standalone scripts introduced.
    - Script count in `config/bin/scripts` reduced by at least 50% with feature parity.
    - All migrated features are documented under “ArchRiot CLI Flags.”

8. Secure Boot (sbctl) wizard — opt-in, gated, and documented (no auto-run)

- Wizard: detect → keygen → sign → pacman hooks → pre-reboot checklist
- Post-reboot continuation validates and restores normal behavior
- Document clear rollback; verify on supported hardware
