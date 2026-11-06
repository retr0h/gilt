# AGENTS.md

Test: `task test` | Before committing: `task fmt` and `task vet`

Read [docs/docs/development.md](docs/docs/development.md) first. It covers
prerequisites, setup, code style, running your changes, testing, branching and
commit messages. [docs/docs/contributing.md](docs/docs/contributing.md) covers
the pull request workflow, and
[docs/docs/architecture.md](docs/docs/architecture.md) the package layout. All
of it applies to agents exactly as it applies to people. This file carries only
what is specific to agents.

## Running tools

Invoke tools through `mise`, not from your path:

```bash
mise exec -- task test
```

`mise` is active in a person's shell and supplies the versions `.mise.toml`
declares. An agent's shell has no activation, so a bare `go` resolves to
whatever is installed globally, usually an older version.

The symptom is a check that fails here and passes in continuous integration, on
a file nobody edited. When that happens, establish which version ran before
treating the failure as real.

## Task tracking

Work is tracked as markdown files in `.tasks/`, in the format
[.tasks/README.md](.tasks/README.md) describes. At the start of a session, check
`.tasks/in-progress/` for ongoing work and `.tasks/backlog/` for what is next.
Move a task file between directories as its status changes, and log the
session's work in `.tasks/sessions/YYYY-MM-DD.md`.

## Commit trailer

When committing via Claude Code, end the message with:

```
🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```
