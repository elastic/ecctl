---
name: Release Checklist
about: Items to be completed in every release.

---

## Version to be released
<!-- Please write the version to be released. -->
<!-- It should follow the following format: vX.Y.Z -->

## Checklist

The following actions must be performed in the order specified here.
For detailed instructions on all of the steps, read the release [documentation](https://github.com/elastic/ecctl/blob/master/developer_docs/RELEASE.md).

- [ ] Update the cloud-sdk-go dependency.
- [ ] Check that the `VERSION` environment variable in the Makefile
- [ ] Generate change log (both release-notes/index.md and notes/*.md)
- [ ] Execute the release.
- [ ] Alert marketing team to update downloads website with new links.
- [ ] Celebrate :tada:
