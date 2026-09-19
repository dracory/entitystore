# AI Memory Bank

Persistent implementation knowledge for AI agents working on this repository.

## How to Use

**Before starting a task:**

- Read `lessons.md` — it records bugs found in code review, non-obvious design
  decisions, and API conventions (e.g., JSON serialization for slices,
  referential-integrity validation, keeping Count filters identical to List
  filters via shared `applyXFilters` helpers).
- Check `AGENTS.md` at the repo root for commands and current conventions.

**While working:**

- If you discover a correction, a pattern that must be followed, or a design
  decision future agents could regress, append a dated entry to `lessons.md`.
- Keep entries short: problem → solution → one-line lesson. Include a minimal
  code snippet only when it clarifies the pattern.
- Record durable conventions in `AGENTS.md` instead; `lessons.md` is for
  incident-derived knowledge.

**Entry format:**

```markdown
### N. Short Title

**Problem:** What went wrong or was ambiguous.

**Solution:** What was done (code snippet if helpful).

**Lesson:** One-sentence rule to apply next time.
```

## Scope

- `lessons.md` — all implementation lessons, grouped by domain/date.
- Do not store secrets, user data, or generated artifacts here.
- Historical design proposals live in `docs/proposals/` and are not maintained
  as current documentation.
