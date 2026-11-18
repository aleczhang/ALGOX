# Subsystem design references

These pages describe current responsibilities, runtime paths, state and
contracts. Start with the subsystem your task touches, then follow its
source links to verify what you change.

| Subsystem | Reference |
| --- | --- |
| Desktop window, the page, its JSON API, Browser access and `algox web` | [GUI app shell](gui-shell.md) |

## Maintain a reference

Update the owning reference in the same change when it alters a documented
responsibility, interface, state transition, runtime path or failure
behavior. A refactor that keeps these contracts needs no prose change.

Keep each fact in one reference and link to it. Link to source files and
named symbols instead of copying details that change often. A new reference
covers responsibilities and boundaries, source of truth, runtime paths,
state, constraints, failure behavior and the tests that verify it.

## Describe a change

For each affected subsystem, say the observable change: its trigger and the
behavior before and after. Link the reference, the code and the test.

## Review a change

1. Find the affected subsystem and the actual runtime path.
2. Compare the claimed before/after with the code, beyond the changed lines.
3. Check the reference still matches and its links resolve.
4. Check the cited tests exercise the claim. Report checks that weren't run.
