# ArchRiot Development Plan — 3.8+

Purpose

- Keep this plan tight, actionable, and future‑proof.
- Capture hard‑won lessons up front, then list only the next steps.
- One change at a time; green builds; explicit commits.

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
- Prefer Go over shell scripts for first‑class features; when shells remain, guard their calls and paths.

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

2. Memory defaults (document and verify)
    - Document conservative defaults and make clear that advanced tuning is opt‑in.
    - Add simple sanity checks around memory‑related features to avoid unintended changes on unsupported systems.

3. Monolith refactors (incremental; one file at a time)
    - Split Waybar emitters from source/waybar/json.go into dedicated files:
        - waybar/pomodoro.go (EmitPomodoro)
        - waybar/memory.go (EmitMemory)
        - waybar/cpu.go (EmitCPU)
        - waybar/temp.go (EmitTemp)
        - waybar/volume.go (EmitVolume)
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

7. Release readiness (3.8.x cadence)
    - VERSION on master must match the badge; the raw VERSION endpoint drives update checks.
    - Prepare a 3.8.1 patch path if regressions appear (e.g., launcher class/routing, help opener, installer guards).

Runtime Validation (Quick Checklist)

- Core commands:
    - ./install/archriot --version
    - ./install/archriot --help
- Waybar lifecycle:
    - ./install/archriot --waybar-status
    - ./install/archriot --waybar-reload
    - ./install/archriot --waybar-sweep
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
