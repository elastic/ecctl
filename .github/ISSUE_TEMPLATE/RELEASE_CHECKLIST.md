---
name: Release Checklist
about: Items to be completed in every release.
labels: Team:Delivery

---

## Version to be released
<!-- Please write the version to be released. -->
<!-- It should follow the following format: vX.Y.Z -->

## Checklist

The following actions must be performed in the order specified here.
For detailed instructions on all of the steps, read the release [documentation](https://github.com/elastic/ecctl/blob/master/developer_docs/RELEASE.md).

- [ ] Verify a cloud-sdk-go release with the same version has been made.
- [ ] Update the cloud-sdk-go dependency.
- [ ] Update `VERSION` environment variable in the Makefile and generate changelog.
- [ ] Execute the release.
- [ ] Modify the conf.yaml file in the docs repository to add the release branch.
- [ ] Run a full doc rebuild.
- [ ] Alert marketing team to update downloads website with new links.
- [ ] Celebrate :tada:
