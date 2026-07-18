---
title: Replace git CLI dependency with go-git
status: in-progress
created: 2026-02-16
updated: 2026-07-18
---

## Objective

Remove the implicit runtime dependency on the git CLI by replacing it with
[go-git](https://github.com/go-git/go-git). This eliminates the requirement
for git >= 2.20 to be installed and in `$PATH`, making Gilt fully
self-contained.

go-git supports bare clones and has enough worktree support for what Gilt
needs (clone, checkout specific version, copy files out).

## Blocker

Gilt uses `--filter=blob:none` for efficient partial clones. go-git does not
yet support this filter type. Tracked upstream at
[go-git/go-git#1381](https://github.com/go-git/go-git/issues/1381) — we are
not tracking that thread; we detect the fix by running our own tests against
`@main`, per the loop below.

Related: [retr0h/gilt#72](https://github.com/retr0h/gilt/issues/72)

## Background

The current git CLI wrapper lives in `internal/git/` and implements the
`GitManager` interface. Because the interface is abstract, a go-git backend
can be swapped in later with minimal changes to the rest of the codebase.

Two files are relevant to this issue:

- `internal/git/git.go`
- `internal/git/git_public_test.go`

## The patch under test

This patch removes the workaround that disabled the `blob:none` partial-clone
filter, and restores the test assertion that expects it. It represents the
"real" fix — it should work once go-git supports the filter, and currently
does not.

This patch is **not committed anywhere**. It exists only as the text below.
Do not leave it applied in the working tree at the end of a session (see
Step 5 of the loop).

```patch
diff --git a/internal/git/git.go b/internal/git/git.go
index afbb9ec..c13c95e 100644
--- a/internal/git/git.go
+++ b/internal/git/git.go
@@ -88,11 +88,6 @@ func (g *Git) Clone(gitURL, origin, cloneDir string) error {
 		NoCheckout: true,
 		Filter:     packp.FilterBlobNone(),
 	}
-	// NOTE(nic): blobless clones don't work quite right yet, so turn them off for
-	//  now.  This regression makes Gilt nigh-unusable, but we can knock all the
-	//  rough edges off of everything else while we wait for upstream to finish
-	//  implementing lazy-fetches
-	opts.Filter = packp.Filter("")
 	repo, err := g.gitClone(cloneDir, opts)
 	g.logger.Debug(
 		"git.Clone",
diff --git a/internal/git/git_public_test.go b/internal/git/git_public_test.go
index 2ecf163..5c615da 100644
--- a/internal/git/git_public_test.go
+++ b/internal/git/git_public_test.go
@@ -182,9 +182,7 @@ func (suite *GitManagerPublicTestSuite) TestClone() {
 	assert.NoError(suite.T(), err)
 	opts := cfg.Raw.Section("remote").Subsection(suite.origin).Options
 	assert.Equal(suite.T(), "true", opts.Get("promisor"))
-	// NOTE(nic): use this assertion instead when blobless clone support is finished
-	// assert.Equal(suite.T(), "blob:none", opts.Get("partialclonefilter"))
-	assert.Equal(suite.T(), "", opts.Get("partialclonefilter"))
+	assert.Equal(suite.T(), "blob:none", opts.Get("partialclonefilter"))
 }
 
 func (suite *GitManagerPublicTestSuite) TestCloneErrorOnOpen() {

```

## Validation loop

Run this once per check-in on upstream progress. Assumes the repo is already
cloned and tooling (Go, `task`) is already installed; `go.mod` is the source
of truth for module versions.

### Step 1 — Bump dependencies to latest `main`

```bash
go get github.com/go-git/go-git/v6@main
go get github.com/go-git/go-billy/v6@main
```

### Step 2 — Test A: confirm the branch still works without the patch

```bash
task test
```

Expected result: **pass**. This is the current shipped behavior (workaround
still active) and should never break. If this fails, something regressed
upstream that is unrelated to issue #72 — stop and investigate separately
before continuing.

### Step 3 — Apply the patch and run Test B: check if upstream has fixed the filter

```bash
git apply .tasks/in-progress/2026-02-16-replace-git-cli-with-go-git.md
task test
```

### Step 4 — Branch on the result of Test B

- **If Test B fails** (expected/known-blocked state, not a regression):
  1. Record a dated entry under "Tracking notes": dependency versions
     tested, result = still blocked.
  2. Revert the patch:
     `git apply -R .tasks/in-progress/2026-02-16-replace-git-cli-with-go-git.md`
  3. Confirm `git status` is clean aside from the kept `go.mod`/`go.sum`
     bump from Step 1 (see "Note on dependency bumps" below).
  4. Session ends here. Re-run this loop next time you want to check
     upstream progress.

- **If Test B passes** (go-git now supports `blob:none` — the loop's exit
  condition has been met):
  1. Record a dated entry under "Tracking notes": dependency versions
     tested, result = **resolved**, this becomes the closing entry.
  2. Do **not** revert the patch. It is no longer speculative — it is the
     real fix.
  3. Remove the "expected to fail" framing from
     `internal/git/git_public_test.go` and delete the old commented-out
     assertion line that this patch replaces (see "The patch under test"
     above — the comment above the removed assertion is now obsolete).
  4. Commit the change as a normal PR against `internal/git/git.go` and
     `internal/git/git_public_test.go`, referencing
     [retr0h/gilt#72](https://github.com/retr0h/gilt/issues/72) and
     [go-git/go-git#1381](https://github.com/go-git/go-git/issues/1381)
     as resolved.
  5. This validation loop is now retired for this issue — do not run it
     again once the PR merges.

## Note on dependency bumps

`go.mod`/`go.sum` changes from Step 1 are **kept, not reverted**, between
sessions. Each run intentionally tracks `go-git`/`go-billy` forward; there is
no rollback step for the dependency bump itself.

## Tracking notes

_Add one line per session: datetime, dependency versions tested, Test B result._

- 2026-07-18: github.com/go-git/go-git/v6 v6.0.0-alpha.4.0.20260716142645-5f90b841aef2, github.com/go-git/go-billy/v6 latest main; Test B still blocked after re-running the validation loop; the patched run failed in integration tests with 14 failures, and the patch was reverted.

## Outcome

Blocked on upstream go-git partial clone support. No further code changes
expected until Test B (Step 3) passes.
