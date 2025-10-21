# ArchRiot Development Plan — Refreshed (v3.5.1)

Purpose

- Keep this plan accurate and actionable.
- Remove completed items from open tasks.
- Capture newly completed work and define focused next steps.
- Maintain a strict one-change-at-a-time cadence: propose → edit → make → verify.

---

Completed in v3.5.1 (shipped)

- README.md overhaul: Navigation/“📚 Navigate This Guide”, Quick Install panel, Advanced Usage/Troubleshooting consolidation, Workspace Styles doc, Contributors section, “What’s New” moved to end
- Keybindings Help PWA: SUPER+SHIFT+H opens local help (Brave app) with stable/fallback window rules; GTK fallback available; binding updated
- Mullvad/VPN stability: Removed NetworkManager reload during install to avoid VPN drop during installation
- Preflight fixes: robust portal detection (PATH, /usr/lib, or running process); cleaned Wi‑Fi power-save diagnostics output
- Installer guards: pre-install ABI advisory and --strict-abi enforcement for compositor/Wayland packages

Next Steps (upgrade-safe; one change at a time)

2. Issue verification and closure (validate, then close)

- #31 Brave crash on workspace/screen switch: retest across GPUs; document safe-mode flag if needed
- #30 OBS “Screen Capture” present: portals functional on clean install
- #24 Thunar default terminal: verify Ghostty “Open Terminal Here” on a clean install

3. README.md finalization (post‑overhaul polish)

- Final pass to remove any duplication and normalize heading styles
- Ensure all code snippets are fenced; fix spacing/indentation regressions
- Keep “What’s New” at end; ensure “Navigate This Guide” is the only TOC
- Confirm install method deep links (Install Script, ISO) present in Navigate/Quick Install

4. QA matrix and docs refresh

- Matrix: Intel/AMD/NVIDIA; single/multi-monitor; fractional scaling; en-US + one non-Latin
- Record pass/fail and remediation steps; update README with confirmed guidance/screenshots

5. HyprToolkit migration (Control Panel first; feature-flagged)

6. GTK Keybindings Help (native app; PWA shipped)

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
