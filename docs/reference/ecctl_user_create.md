---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_user_create.html
---

# ecctl user create [ecctl_user_create]

Creates a new platform user ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl user create --username <username> --role <role> [flags]
```


## Examples [_examples_13]

```
  * Create a platform user who has two roles assigned
    ecctl user create --username sam89 --role ece_platform_viewer --role ece_deployment_viewer
```


## Options [_options_126]

```
      --email string               User's email address (must be in a valid email format)
      --fullname string            User's full name
  -h, --help                       help for create
      --insecure-password string   [INSECURE] User's plaintext password
      --role strings               Role or roles assigned to the user. Available roles:
                                   ece_platform_admin, ece_platform_viewer, ece_deployment_manager, ece_deployment_viewer
      --username string            Unique username for the platform user
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_125]

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_126]

* [ecctl user](/reference/ecctl_user.md)	 - Manages the platform users ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

