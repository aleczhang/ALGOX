# Lessons

What merged work got wrong, and the rule that would have caught it. Newest
first. Each entry: the date, what happened (with the commit or issue), and
the rule.

<!-- Example:
## 2025-01-01: a setting saved from the browser bypassed its check
`POST /api/settings` wrote the body without `Settings.Check` (abc1234).
Rule: every write route checks through the package that owns the file.
-->
