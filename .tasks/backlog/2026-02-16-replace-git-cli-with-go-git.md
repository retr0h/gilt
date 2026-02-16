---
title: Replace git CLI dependency with go-git
status: backlog
created: 2026-02-16
updated: 2026-02-16
---

## Objective

Remove the implicit runtime dependency on the git CLI by replacing it with
[go-git](https://github.com/go-git/go-git). This would eliminate the
requirement for git >= 2.20 to be installed and in `$PATH`, making Gilt fully
self-contained.

go-git supports bare clones and has enough worktree support for what Gilt needs
(clone, checkout specific version, copy files out).

## Blockers

- **Blocked on go-git partial clone support:**
  [go-git/go-git#1381](https://github.com/go-git/go-git/issues/1381).
  Gilt uses `--filter=blob:none` for efficient partial clones; go-git does not
  yet support this.

## Notes

- Related issue: [retr0h/gilt#72](https://github.com/retr0h/gilt/issues/72)
- The current git CLI wrapper lives in `internal/git/` and implements the
  `GitManager` interface. The interface-based design means go-git can be
  swapped in as an alternative implementation with minimal changes to the
  rest of the codebase.
- Consider a phased approach: implement go-git backend behind a flag first,
  fall back to CLI for partial clone operations until go-git catches up.

## Outcome

_To be filled in when complete._
