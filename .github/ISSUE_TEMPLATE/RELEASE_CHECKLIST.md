---
name: Release Checklist
about: Items to be completed in every release.
labels: area:delivery

---
<!-- Before starting the release process please make sure you've -->
<!-- read the release documentation in /developer_docs/RELEASE.md -->

## Version to be released
<!-- Please write the version to be released. -->
<!-- It should follow the following format: vX.Y.Z -->

## Checklist
<!-- The following actions must be performed in the order specified here. -->
- [ ] Verify a cloud-sdk-go release with the same version has been made.
- [ ] Update the cloud-sdk-go dependency.
- [ ] Update `VERSION` environment variable in the Makefile and generate changelog.
- [ ] Execute the release.
- [ ] Create release branch for documentation.
- [ ] Modify the conf.yaml file in the docs repository to add the release branch.
- [ ] Run a full doc rebuild.
- [ ] Celebrate :tada:
