---
sidebar_position: 8
---

# Architecture

**Date:** February 2026 **Author:** @retr0h

## Overview

Gilt is a repository overlay tool. It clones Git repositories at pinned versions
(tags or SHAs), then copies specific files or directories into a target project.
This enables vendoring configuration, schemas, or shared code from upstream
repositories without Git submodules.

## High-Level Flow

```
Giltfile.yaml
     │
     v
┌─────────────────────────┐
│  pkg/repositories       │  Public API entry point
│  Overlay()              │
└────────────┬────────────┘
             │
             v
┌─────────────────────────┐
│  internal/repositories  │  Iterates configured repos (serial or parallel)
│  Overlay()              │
└────────────┬────────────┘
             │
             v
┌─────────────────────────┐
│  internal/repository    │  Per-repo: clone, worktree, copy sources
│  Clone() / Worktree()   │
│  CopySources()          │
└────────────┬────────────┘
             │
             v
┌─────────────────────────┐
│  internal/git           │  Shells out to git CLI
│  Clone() / Worktree()   │  (bare clone, worktree checkout)
│  Update() / Remote()    │
└────────────┬────────────┘
             │
             v
┌─────────────────────────┐
│  internal/exec          │  Command execution abstraction
│  RunCmd()               │
│  RunCmdInDir()          │
└─────────────────────────┘
```

## Package Details

### `cmd/`

Cobra CLI commands. Entry points for user interaction.

- `root.go` - Root command, Viper config binding
- `overlay.go` - `gilt overlay` command, reads Giltfile and runs the overlay
- `init.go` - `gilt init` command, scaffolds a new Giltfile
- `version.go` - `gilt version` command

### `internal/`

Interface definitions live at the package root (`git.go`, `exec.go`,
`repository.go`, `repositories.go`). Implementations live in sub-packages.

- **`git/`** - Git CLI wrapper. Performs bare clones, worktree checkouts, remote
  URL lookups, and repository updates by shelling out to `git`.
- **`exec/`** - Command execution abstraction. Wraps `os/exec` with working
  directory support and temp directory helpers.
- **`repository/`** - Single repository operations. Orchestrates clone, worktree
  checkout, and file/directory copying for one repository entry.
- **`repositories/`** - Multi-repository orchestrator. Reads the Giltfile,
  iterates all configured repositories, and delegates to `repository/`. Supports
  parallel execution.
- **`path/`** - Path utility functions.
- **`mocks/`** - Generated mock implementations (via `mockgen`) for all
  interfaces. Used in unit tests.

### `pkg/`

Public API surface.

- **`config/`** - Configuration types (`Repositories`, `Repository`, `Source`,
  `Command`). Uses Viper for binding and `go-playground/validator` for schema
  validation.
- **`repositories/`** - Public entry point. Wires together internal components
  and exposes `Overlay()` to external consumers.

### `test/integration/`

Bats integration tests that exercise the full `gilt overlay` flow against real
Git repositories.

### `docs/`

Docusaurus documentation site (Markdown content in `docs/docs/`).

### `python/`

Python wheel packaging scripts for distributing Gilt via PyPI.

## Design Patterns

### Interface Segregation

Small, focused interfaces are defined in `internal/*.go`:

- `GitManager` - Git operations (clone, worktree, update, remote)
- `ExecManager` - Command execution (run, run-in-dir, run-in-temp-dir)
- `RepositoryManager` - Single repo operations (clone, worktree, copy)
- `RepositoriesManager` - Multi-repo orchestration (overlay)

### Dependency Injection

Implementations accept their dependencies via constructors. For example,
`git.New()` takes an `ExecManager`, and `repository.New()` takes both a
`GitManager` and `ExecManager`. This makes testing straightforward with
generated mocks.

### Filesystem Abstraction

The project uses [`avfs`](https://github.com/avfs/avfs) for filesystem
operations, enabling tests to run against in-memory filesystems.

### Configuration Layering

Viper binds CLI flags, environment variables, and the Giltfile (YAML) into typed
`config.Repositories` structs. Validation is handled by
`go-playground/validator` struct tags.
