# ArchRiot Development Plan — 3.9+

Purpose

- Keep this plan tight, actionable, and future‑proof.
- Capture hard‑won lessons up front, then list only the next steps.
- One change at a time; green builds; explicit commits.

Non‑Negotiable Rules

- One change at a time → make → verify green → stage exact paths → proceed.
- Do not commit unless explicitly told "commit".
- Never add new .sh scripts for first‑class/runtime features; expose Go flags in the archriot binary instead.
- README edits must be appended at the bottom (release notes/changelog only).
- Guard systemd/user unit enable/start with unit existence; never fail installer if a unit is missing.
- Avoid full Hyprland reloads in installer flows; prefer runtime keyword updates or --stabilize-session.
- No environment variables for behavior; always expose flags.
- main.go remains delegation-only; move logic into packages.
- Never modify ~/.config/hypr/monitors.conf while Hyprland is active; skip backup/edits when `hyprctl -j monitors` succeeds; never trigger compositor reloads as part of upgrades.

Lessons Learned

- Hyprland reloads can crash Brave/reset monitors; avoid full reload for window rules.
- kanshi.service may be missing; gate enablement by unit existence and make it non-fatal.
- Stabilize via --stabilize-session; do not invent new shell helpers.
- Do not modify monitors.conf during a live session; gate changes on Hyprland inactivity; prefer kanshi profiles for hotplug; never force compositor reloads from installer paths.

Lessons & Patterns (Keep These Front‑of‑Mind)

- Build discipline:
    - After every change: run make; the build must be green before you proceed.
    - Stage exact paths; confirm with git status -s before doing anything else.
    - Do not commit unless explicitly told “commit”.
- Entry point purity:
    - main.go is delegation‑only: parse flags and call package functions.
    - No feature logic in main.go; if new logic appears, extract to a package first.
- UX correctness and resilience:
    - Prefer native integrations (e.g., Waybar emitters) over external scripts.
    - Provide safe fallbacks (e.g., opener fallbacks for help windows).
    - Guard filesystem operations (chmod/rm) behind directory/file presence checks.
- Defend against technical debt:
    - Remove dead/unreachable code immediately.
    - Do not add env vars for behavior; always expose flags.
    - Consolidate settings; avoid duplication across modules and configs.
    - Never add new .sh scripts for features that can be implemented as Go flags; implement them in the archriot binary and phase out legacy scripts.
- Rebase safety:
    - Expect merges to change configs (e.g., Waybar). Re‑verify critical binds/UX after rebase.
    - Re‑run validation checklists (below) before tagging a release.

Process (Non‑Negotiable)

- One change at a time → make → verify green → stage exact paths → proceed.
- Commit only when explicitly told “commit”.
- Keep diffs minimal and scoped; avoid drive‑by edits.
- Never modify the ISO reference in README.md during normal changes (only during release prep if needed).

Architecture Guardrails

- Flags map 1:1 to package functions (no feature logic in main.go).
- Extract multi‑concern logic to dedicated packages.
- Large modules should be split by concern (emitters, tools, TUI model/view/update, etc.).
- Never add new shell scripts for first‑class/runtime features; expose a Go flag in the archriot binary instead. If legacy shells remain, guard their calls/paths and plan removal.

Quality Gates (Always Validate)

- Waybar routing: native on‑click "activate", workspace persistence (1–4) with 5–6 dynamic.
- Help binder: SUPER+SHIFT+H opens generated Keybindings Help.
- Telegram launcher UX: PATH‑resolved + 2s notify (where applicable).
- Pomodoro transitions: toast notifications; correct status JSON.
- Installer: chmod/rm guards for removed script directories; no hard failures if absent.

Next Steps (In Order)

1. Docs & QA polish
    - README alignment with current behavior:
        - Workspace persistence (1–4), dynamic 5–6.
        - Native “activate” for workspace clicks (per‑monitor routing).
        - SUPER+SHIFT+H help behavior; Telegram minimal launcher UX.
    - Add a concise QA checklist (see Runtime Validation) and a short “Contributing” blurb referencing this plan’s rules.

2. Memory stability (verify and document)
    - Verify zram (zstd, size = ram/2) and sysctl tuning are applied; document opt‑out/override in README.
    - Ensure systemd‑oomd remains disabled/masked during install/upgrade.

3. Monolith refactors (incremental; one file at a time)
    - Split tools/tools.go by tool (secure boot, memory optimizer, perf tuner, dev env) and keep a tiny registry/factory.
    - Split TUI files: model.go, update.go, view.go, messages.go for clarity.

4. Installer safeguards (audit for durability)
    - Audit install/packages.yaml for lingering references to removed directories or script globs.
    - Guard all chmod/rm/ln calls with presence checks; always ensure these are non‑fatal when targets are missing.

5. Feature‑gated launcher evaluation
    - Add Hyprlauncher behind a feature flag; keep Fuzzel as default.
    - Validate parity and stability (focus behavior, classes, warm/cold start).
    - Only consider flipping the default when parity is perfect across configs.

6. Tests & traceability (targeted)
    - Help system: add unit tests for key normalization/display for binds; minimal test harness for HTML generation (no headless browser).
    - Help opener: optional debug mode to emit the chosen opener and path when troubleshooting (silent by default).

7. Release readiness (3.9.x cadence)
    - VERSION on master must match the badge; the raw VERSION endpoint drives update checks.
    - Prepare a 3.9.x patch path if regressions appear.

Runtime Validation (Quick Checklist)

- Core commands:
    - ./install/archriot --version
    - ./install/archriot --help
- Waybar & Displays lifecycle:
    - ./install/archriot --waybar-status
    - ./install/archriot --waybar-reload
    - ./install/archriot --waybar-sweep
    - ./install/archriot --displays-enforce (verify external‑only when docked; laptop‑only when undocked)
    - systemctl --user is-enabled --quiet archriot-displays-enforce.service (enabled at login)
    - systemctl --user is-active --quiet kanshi.service (running; may be enabled‑only if Brave gating applied)
- Workspaces:
    - Validate persistent 1–4 only; 5–6 appear when in use.
    - Clicks route correctly on the clicked monitor (native “activate”).
- Help system:
    - SUPER+SHIFT+H opens Keybindings Help (Brave app if present; gtk‑launch/xdg‑open fallback); generated path exists.
- Pomodoro:
    - Toggle & reset behavior; toast notifications on work/break transitions.
- Installer robustness:
    - No hard failures if ~/.local/share/archriot/config/bin/scripts or ~/.config/waybar/scripts are missing.

Git: safe pull with local changes

- Preferred (one-shot):
    - git pull --rebase --autostash

- Manual stash, then pull:
    - git stash push -u -m "WIP before pull"
    - git pull --rebase
    - git stash pop

- Isolate WIP (optional, safest for bigger changes):
    - git switch -c wip/keep-local
    - git switch master && git pull --rebase
    - git switch wip/keep-local && git rebase master

Always:

- Check git status -s before/after
- If conflicts occur: resolve, then git add -A && git rebase --continue

Appendix: Quick Commands

- Build: make
- Version: ./install/archriot --version
- Help binder: ./install/archriot --help-binds-generate and ./install/archriot --help-binds-web
- Waybar reload: ./install/archriot --waybar-reload
- Workspace click (for custom modules): ./install/archriot --waybar-workspace-click {name}
- Kanshi autogen: ./install/archriot --kanshi-autogen
