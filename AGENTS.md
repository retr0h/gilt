# AGENTS.md

This file provides guidance to AI coding agents when working with code in this
repository.

## Project Overview

Gilt is a Git repository overlay tool written in Go 1.25. It clones Git
repositories, checks out specific versions (tags or SHAs), and copies files or
directories into a target project. Module path: `github.com/retr0h/gilt/v2`.

## Development Reference

For setup, building, testing, and contributing, see the Docusaurus docs:

- @docs/docs/development.md - Prerequisites, setup, code style, testing, commit
  conventions
- @docs/docs/contributing.md - PR workflow and contribution guidelines
- @docs/docs/testing.md - How to run tests and list task recipes
- @docs/docs/architecture.md - Package architecture and design patterns

Quick reference for common commands:

```bash
task deps          # Install all dependencies
task test          # Run all tests (lint + unit + coverage + bats)
task unit          # Run unit tests only
task vet           # Run golangci-lint
task fmt           # Auto-format (gofumpt + golines)
task fmt:check     # Check formatting without modifying
go test -run TestName -v ./internal/...  # Run a single test
```

## Architecture (Quick Reference)

- **`cmd/`** - Cobra CLI commands (`root`, `overlay`, `init`, `version`)
- **`internal/`** - Interface definitions (`git.go`, `exec.go`, `repository.go`,
  `repositories.go`)
- **`internal/git/`** - Git CLI wrapper (clone, worktree, update, remote)
- **`internal/exec/`** - Command execution abstraction
- **`internal/repository/`** - Single repository operations (clone, worktree,
  copy sources)
- **`internal/repositories/`** - Orchestrates overlay across all configured
  repositories
- **`internal/path/`** - Path utility functions
- **`internal/mocks/`** - Generated mocks (mockgen)
- **`pkg/config/`** - Configuration types (`Repositories`, `Repository`,
  `Source`, `Command`) with Viper + validator tags
- **`pkg/repositories/`** - Public API entry point for repository operations
- **`test/integration/`** - Bats integration tests
- **`docs/`** - Docusaurus documentation site
- **`python/`** - Python wheel packaging for PyPI distribution

## Code Standards (MANDATORY)

### Testing

- Unit tests: `*_test.go` in same package for private functions
- Public tests: `*_public_test.go` in test package (e.g., `package git_test`)
  for exported functions
- Integration tests in `test/integration/` using Bats
- Use `testify/assert` and `testify/require`

### Go Patterns

- Interface segregation: small interfaces in `internal/*.go`, implementations
  in sub-packages
- Dependency injection via constructors (e.g., `NewGit(execManager)`)
- Error wrapping: `fmt.Errorf("context: %w", err)`
- Early returns over nested if-else
- Unused parameters: rename to `_`
- Import order: stdlib, third-party, local (blank-line separated)

### Linting

golangci-lint with: errcheck, errname, goimports, govet, prealloc, predeclared,
revive, staticcheck. Formatting via gofumpt + golines.

### Branching

See @docs/docs/development.md#branching for full conventions.

When committing changes, create a feature branch first if currently on `main`.
Branch names use the pattern `type/short-description` (e.g.,
`feat/add-dns-retry`, `fix/memory-leak`, `docs/update-readme`).

### Commit Messages

See @docs/docs/development.md#commit-messages for full conventions.

Follow [Conventional Commits](https://www.conventionalcommits.org/) with the
50/72 rule. Format: `type(scope): description`.

## Task Tracking

Work is tracked as markdown files in `.tasks/`. See @.tasks/README.md for
format details.

```
.tasks/
├── backlog/          # Tasks not yet started
├── in-progress/      # Tasks actively being worked on
├── done/             # Completed tasks
└── sessions/         # Session work logs (per session)
```

When starting a session:

1. Check `.tasks/in-progress/` for ongoing work
2. Check `.tasks/backlog/` for next tasks
3. Move task files between directories as status changes
4. Log session work in `.tasks/sessions/YYYY-MM-DD.md`
