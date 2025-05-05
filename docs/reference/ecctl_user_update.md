---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_user_update.html
---

# ecctl user update [ecctl_user_update]

Updates a platform user ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl user update <username> --role <role> [flags]
```


## Examples [_examples_14]

```
  * Update a platform user
	ecctl user update --username xochitl --role ece_platform_viewer --email xo@example.com

  * Update the current platform user
    ecctl user update --current --email xo@example.com --password
```


## Options [_options_136]

```
      --current                    updates details of the current user
      --email string               user's email address (must be in a valid email format)
      --fullname string            user's full name
  -h, --help                       help for update
      --insecure-password string   [INSECURE] user plaintext password
      --password                   if set, updates the user's password securely
                                   (must use a minimum of 8 characters )
      --role strings               role or roles assigned to the user. Available roles:
                                   ece_platform_admin, ece_platform_viewer, ece_deployment_manager, ece_deployment_viewer
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_135]

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_136]

* [ecctl user](/reference/ecctl_user.md)	 - Manages the platform users ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

