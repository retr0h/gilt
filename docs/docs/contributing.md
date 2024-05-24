---
sidebar_position: 6
---

# Contributing

Contributions to Gilt are very welcome, but we ask that you read this document
before submitting a PR.

:::note

This document applies to the [Gilt][] repository.

:::

## Before you start

- **Check existing work** - Is there an existing PR? Are there issues discussing
  the feature/change you want to make? Please make sure you consider/address
  these discussions in your work.
- **Backwards compatibility** - Will your change break existing Giltfiles? It is
  much more likely that your change will merged if it backwards compatible. Is
  there an approach you can take that maintains this compatibility? If not,
  consider opening an issue first so that API changes can be discussed before
  you invest your time into a PR.

## 1. Setup

- **Go** - Gilt is written in [Go][]. We always support the latest two major Go
  versions, so make sure your version is recent enough.
- **Node.js** - [Node.js][] is used to host Gilt's documentation server and is
  required if you want to run this server locally.

## 2. Making changes

- **Code style** - Try to maintain the existing code style where possible. Go
  code should be formatted by [`gofumpt`][gofumpt] and linted using
  [`golangci-lint`][golangci-lint]. Any Markdown or TypeScript files should be
  formatted and linted by [Prettier][]. This style is enforced by our CI to
  ensure that we have a consistent style across the project. You can use the
  `task fmt:check` command to lint the code locally and the `task fmt` command
  to automatically fix any issues that are found.
- **Documentation** - Ensure that you add/update any relevant documentation. See
  the [updating documentation](#updating-documentation) section below.
- **Tests** - Ensure that you add/update any relevant tests and that all tests
  are passing before submitting the PR. See the [writing tests](#writing-tests)
  section below.

### Running your changes

To run Gilt with working changes, you can use `go run main.go overlay`.

### Updating documentation

Gilt uses [Docusaurus][] to host a documentation server. The code for this is
located in the Gilt repository. This can be setup and run locally by using
`task docs:start` (requires `nodejs` & `yarn`). All content is written in
Markdown and is located in the `docs/docs` directory. All Markdown documents
should have an 80 character line wrap limit (enforced by Prettier).

### Writing tests

When making changes, consider whether new tests are required. These tests should
ensure that the functionality you are adding will continue to work in the
future. Existing tests may also need updating if you have changed Gilt's
behavior.

You may also consider adding unit tests for any new functions you have added.
The unit tests should follow the Go convention of being location in a file named
`*_test.go` in the same package as the code being tested.

Integration tests are located in the `tests` directory and executed by [Bats][].

## 3. Committing your code

Try to write meaningful commit messages and avoid having too many commits on the
PR. Most PRs should likely have a single commit (although for bigger PRs it may
be reasonable to split it in a few). Git squash and rebase is your friend!

If you're not sure how to format your commit message, check out [Conventional
Commits][]. This style is enforced, and is a good way to make your commit
messages more readable and consistent.

## 4. Submitting a PR

- **Describe your changes** - Ensure that you provide a comprehensive
  description of your changes.
- **Issue/PR links** - Link any previous work such as related issues or PRs.
  Please describe how your changes differ to/extend this work.
- **Examples** - Add any examples or screenshots that you think are useful to
  demonstrate the effect of your changes.
- **Draft PRs** - If your changes are incomplete, but you would like to discuss
  them, open the PR as a draft and add a comment to start a discussion. Using
  comments rather than the PR description allows the description to be updated
  later while preserving any discussions.

## FAQ

> I want to contribute, where do I start?

All kinds of contributions are welcome, whether its a typo fix or a shiny new
feature. You can also contribute by upvoting/commenting on issues or helping to
answer questions.

> I'm stuck, where can I get help?

If you have questions, feel free open a [Discussion][] on GitHub.

<!-- prettier-ignore-start -->
[Gilt]: https://github.com/retr0h/gilt
[Go]: https://go.dev
[Node.js]: https://nodejs.org/en/
[gofumpt]: https://github.com/mvdan/gofumpt
[golangci-lint]: https://golangci-lint.run
[Prettier]: https://prettier.io/
[Docusaurus]: https://docusaurus.io
[Discussion]: https://github.com/retr0h/gilt/discussions
[Conventional Commits]: https://www.conventionalcommits.org
[Bats]: https://github.com/bats-core/bats-core
<!-- prettier-ignore-end -->
