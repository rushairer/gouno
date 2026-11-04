# Release Checklist

- [ ] Public API and minimum Go version changes are documented in English and Chinese.
- [ ] `CHANGELOG.md` contains the dated v1.4.0 SemVer entry.
- [ ] Go 1.26.x and 1.27.x CI jobs pass on the exact release candidate commit.
- [ ] Formatting, tidy-without-diff, race tests, vet, lint, and govulncheck pass on both supported Go series.
- [ ] No reachable vulnerability or unreviewed dependency source remains.
- [ ] Confirm the v1.4.0 tag does not exist and will point at the reviewed release commit; never reuse a tag.
- [ ] After publishing, verify the public module proxy resolves `github.com/rushairer/gouno@v1.4.0`.
- [ ] After the Core tag is public, update the official template, Gosso, and Blog backend to the exact v1.4.0 dependency where applicable; run tidy and their complete CI.
- [ ] Publish the default Template and CLI in the documented dependency order only after their integration gates pass.
- [ ] Update the bilingual compatibility matrix only after the corresponding release tags and packages are publicly available.
