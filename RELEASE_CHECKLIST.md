# Release Checklist

- [ ] Public API and minimum Go version changes are documented.
- [ ] `CHANGELOG.md` has a dated SemVer entry.
- [ ] Go 1.25.x and 1.26.x CI jobs pass.
- [ ] Tests, race detector, vet, lint, and govulncheck pass.
- [ ] No reachable vulnerability or unreviewed dependency source remains.
- [ ] Tag `vX.Y.Z` points at the reviewed commit and has not been reused.
