---
sidebar_position: 7
---

# Development

This guide covers the tools, setup, and conventions needed to work on Gilt.

## Prerequisites

Install tools using [mise][]:

```bash
mise install
```

- **[Go][]** - Gilt is written in Go. We always support the latest two major Go
  versions, so make sure your version is recent enough.
- **[Node.js][]** - Required for the Docusaurus documentation site.
- **[go-task][]** - Task runner used for building, testing, formatting, and
  other development workflows.

### Claude Code

If you use [Claude Code][] for development, install the **commit-commands**
plugin from the default marketplace:

```
/plugin install commit-commands@claude-plugins-official
```

This provides `/commit` and `/commit-push-pr` slash commands that follow the
project's commit conventions automatically.

## Setup

Install dependencies and verify your environment:

```bash
task deps          # Install Go tool dependencies
task deps:check    # Verify required tools are available
```

## Code style

Go code should be formatted by [`gofumpt`][gofumpt] and [`golines`][golines],
and linted using [`golangci-lint`][golangci-lint]. Markdown and TypeScript files
should be formatted and linted by [Prettier][]. This style is enforced by CI.

```bash
task fmt:check     # Check formatting
task fmt           # Auto-fix formatting
task vet           # Run linter
```

## Running your changes

To run Gilt with working changes:

```bash
task run -- overlay
```

Or directly:

```bash
go run main.go overlay
```

## Documentation

Gilt uses [Docusaurus][] to host a documentation server. Content is written in
Markdown and located in the `docs/docs` directory. All Markdown documents should
have an 80 character line wrap limit (enforced by Prettier).

```bash
task docs:start    # Start local docs server (requires bun)
```

## Testing

See the [Testing](testing.md) page for details on running tests.

```bash
task test          # Run all tests (lint + unit + coverage + bats)
task unit          # Run unit tests only
task unit:int      # Run integration tests only (Bats)
```

Unit tests should follow the Go convention of being located in a file named
`*_test.go` in the same package as the code being tested. Public API tests use
the `*_public_test.go` convention with an external test package (e.g.,
`package git_test`). Integration tests are located in the `test/integration`
directory and executed by [Bats][].

## Branching

All changes should be developed on feature branches. Create a branch from `main`
using the naming convention `type/short-description`, where `type` matches the
[Conventional Commits][] type:

- `feat/add-retry-logic`
- `fix/null-pointer-crash`
- `docs/update-api-reference`
- `refactor/simplify-handler`
- `chore/update-dependencies`
- `test/add-clone-tests`
- `style/fix-formatting`
- `perf/optimize-copy`

When using Claude Code's `/commit` command, a branch will be created
automatically if you are on `main`.

## Commit messages

Follow [Conventional Commits][] with the 50/72 rule:

- **Subject line**: max 50 characters, imperative mood, capitalized, no period
- **Body**: wrap at 72 characters, separated from subject by a blank line
- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- Summarize the "what" and "why", not the "how"

Try to write meaningful commit messages and avoid having too many commits on a
PR. Most PRs should likely have a single commit (although for bigger PRs it may
be reasonable to split it in a few). Git squash and rebase is your friend!

<!-- prettier-ignore-start -->
[mise]: https://mise.jdx.dev
[Go]: https://go.dev
[Node.js]: https://nodejs.org/en/
[go-task]: https://taskfile.dev
[Claude Code]: https://claude.ai/code
[gofumpt]: https://github.com/mvdan/gofumpt
[golines]: https://github.com/segmentio/golines
[golangci-lint]: https://golangci-lint.run
[Prettier]: https://prettier.io/
[Docusaurus]: https://docusaurus.io
[Conventional Commits]: https://www.conventionalcommits.org
[Bats]: https://github.com/bats-core/bats-core
<!-- prettier-ignore-end -->
