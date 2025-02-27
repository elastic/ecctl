---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-example-delete-deployment.html
---

# Delete a deployment [ecctl-example-delete-deployment]

In this last example, you can use the [ecctl deployment shutdown](/reference/ecctl_deployment_shutdown.md) command to delete the deployment that you created.

```sh
ecctl deployment shutdown [--track] $DEPLOYMENT_ID
```

* `$DEPLOYMENT_ID` is the ID for the deployment that was created in the previous [create a deployment](/reference/ecctl-example-create-deployment.md) example.

On running this and other destructive commands, ecctl prompts you with a confirmation message. Use the `--force` option to skip the confirmation step, if you are using ecctl for automation.

To monitor the progress, use the `--track` flag.

To see the different options that ecctl supports, run `ecctl <command> <help>`.

