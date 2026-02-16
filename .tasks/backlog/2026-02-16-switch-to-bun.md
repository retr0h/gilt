---
title: Switch docs site package manager to bun
status: backlog
created: 2026-02-16
updated: 2026-02-16
---

## Objective

Replace yarn/npm with bun as the package manager for the Docusaurus
documentation site.

## Notes

- `.mise.toml` already includes `bun = "latest"` so the runtime is available.
- Requires updating `Taskfile.yml` docs tasks (currently reference yarn).
- Update `package.json` scripts if needed.
- Update CI workflows to use bun instead of yarn/npm.
- Update `docs/docs/contributing.md` references to yarn.

## Outcome

_To be filled in when complete._
