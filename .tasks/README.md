# Task Tracking

Work is tracked as markdown files organized by status.

## Directory Structure

```
.tasks/
├── backlog/          # Tasks not yet started
├── in-progress/      # Tasks actively being worked on
├── done/             # Completed tasks
└── sessions/         # Session work logs (what was done per Claude Code session)
```

## Task File Format

Filename: `YYYY-MM-DD-short-description.md`

```markdown
---
title: Short description of the task
status: backlog | in-progress | done
created: YYYY-MM-DD
updated: YYYY-MM-DD
---

## Objective

What needs to be done and why.

## Notes

Decisions, context, blockers, or references.

## Outcome

What was accomplished (filled in when done).
```

## Session Log Format

Filename: `YYYY-MM-DD.md`

```markdown
---
date: YYYY-MM-DD
---

## Summary

Brief overview of what was accomplished.

## Changes

- List of files changed and why.

## Decisions

- Any architectural or design decisions made.

## Next Steps

- What to pick up next session.
```

## Workflow

1. New work goes in `backlog/`
2. Move file to `in-progress/` when starting
3. Move file to `done/` when complete
4. Log each Claude Code session in `sessions/`
