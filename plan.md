# ArchRiot Development Plan — Next Steps

This plan lists only what remains. All completed items have been removed.

Non‑negotiables (keep us disciplined)

- One change at a time → make → verify green → stage exact paths → proceed (no commit unless explicitly told “commit”).
- main.go is delegation-only (parse flags; wire to packages). No feature logic belongs there.
- No new env vars for behavior (always use flags). No duplicate code. Remove dead/unreachable code immediately.

Priority action queue (in order)

1. Docs & QA polish

- README: ensure all examples match current behavior (workspace 1–4 persistence with dynamic 5–6, native on-click “activate”, SUPER+G Telegram with 2s notify).
- Add a concise QA checklist (install/upgrade, multi-monitor routing, Pomodoro transitions, SUPER+SHIFT+H help).
- Plan a short “Contributing” section pointing to these rules.

2. Memory defaults (document and verify)

- Document safe, conservative defaults; keep advanced tuning strictly opt‑in.
- Audit current defaults for surprises; add simple sanity checks (no-op if unsupported).

3. Monolith refactors (incremental; one file at a time)

- waybar/json.go (split emitters into files):
    - waybar/pomodoro.go (EmitPomodoro)
    - waybar/memory.go (EmitMemory)
    - waybar/cpu.go (EmitCPU)
    - waybar/temp.go (EmitTemp)
    - waybar/volume.go (EmitVolume)
- tools/tools.go (split each tool into its own file; keep a small registry/factory).
- tui/model.go (split into model.go, update.go, view.go, messages.go for clarity).

4. Installer safeguards

- Audit install/packages.yaml for any lingering references to removed script dirs or globs; guard with [ -d … ] checks (we fixed scripts; repeat for any similar patterns).
- Keep chmod/rm calls robust when targets are missing (|| true and directory existence checks).

5. Feature‑gated launcher evaluation

- Add Hyprlauncher behind a feature flag (default remains Fuzzel).
- Validate parity and stability (cold/warm launch, class names, focus behavior). Only consider switching default when parity is perfect.

6. Tests and traceability (targeted)

- Help system: add small unit tests around key normalization/display (bind parsing).
- OpenWeb: add minimal debug logging toggle (opener selected, path existence) for easier field diagnosis (keep silent by default).

7. Release readiness (patch cadence)

- Ensure VERSION only changes when releasing and that the upgrade checker (raw VERSION on master) matches.
- Prepare a 3.7.1 patch checklist if any regressions appear (Waybar routing, help binder, Telegram UX).

Guardrails (do not skip)

- Always run make after each change; the build must be green before staging the next task.
- Stage exact files; verify with git status -s before proceeding.
- If you are not 100% confident a change won’t introduce side effects, do not propose it yet—evaluate and test first.
