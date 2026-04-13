# Tree Challenge

## Context

This package builds a tree of **locations** from a flat input list. A location
can be a Property, a Floor, a Section, a Room, a Storage area, etc. Each
location has an `ID`, a `Name`, an optional `ParentID`, and a `LocationType`.
Exactly one location (the Property) sits at the root.

The `Traverse` function (in `traverse.go`) walks the tree and invokes a
callback once per node. Other parts of the system rely on it — for example,
the assignment pipeline uses `Traverse` to iterate through a property's
locations and schedule work against each one.

## The report

A customer running "Grand Plaza Hotel" reported that when their location tree
is walked with `Traverse`, **some locations are silently skipped**: the
callback never fires for them, so no downstream work gets scheduled against
those locations. No errors, no logs — the nodes are simply missing from the
results.

We reproduced the issue with their data in `traverse_test.go`. That test
currently fails.

## Your task

1. Run the tests and observe the failure.

2. Diagnose the root cause. Be ready to explain *why* the traversal is
   dropping a node.

3. Propose and implement a fix. All tests must pass after your change.

4. Briefly discuss: are there edge cases your fix does **not** cover?

You can run the tests at any point with:

```bash
go test ./...
```

## What we are looking for

- Ability to read unfamiliar code and reason about an observed bug.
- A minimal, well-scoped fix — not a rewrite.
- A clear explanation of root cause and impact.

## Files

- `tree.go` — `Node`, `Location`, `BuildTree`, `GetID`.
- `traverse.go` — `Traverse`.
- `tree_test.go` — sanity test for `BuildTree`.
- `traverse_test.go` — reproduction of the customer report. **Fails today.**
