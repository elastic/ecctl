---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-example-update-deployment.html
---

# Update a deployment [ecctl-example-update-deployment]

Now that you have used ecctl to create a deployment, you can scale it up, by increasing the size of the Elasticsearch data nodes from 1024 to 4096 MB.

Copy the following JSON payload and save it as file `update-deployment.json`.

```json
{
  "prune_orphans": false,
  "resources": {
    "elasticsearch": [
      {
        "region": "gcp-us-central1",
        "ref_id": "main-elasticsearch",
        "plan": {
          "cluster_topology": [
            {
              "node_type": {
                "master": true,
                "data": true,
                "ingest": true,
                "ml": false
              },
              "instance_configuration_id": "gcp.data.highio.1",
              "zone_count": 2,
              "size": {
                "resource": "memory",
                "value": 4096
              }
            }
          ],
          "elasticsearch": {
            "version": "7.6.0"
          },
          "deployment_template": {
            "id": "gcp-io-optimized"
          }
        }
      }
    ]
  }
}
```

The JSON body is similar to what we used to create the deployment, with the following differences:

* The name of the deployment can be modified or it will stay the same if not specified.
* A `prune_orphans` parameter is added. This important parameter specifies how resources not included in the JSON should be handled:

    * if `true`, those resources not included are removed
    * if `false`, those resources not included are kept intact


In this example, prune_orphans is set to `false`, so the Kibana and APM instances are not changed or removed, while the Elasticsearch resource is modified according to the configuration specified in the JSON file.

To monitor the progress, use the `--track` flag.

```sh
ecctl deployment update [--track] $DEPLOYMENT_ID -f update-deployment.json
```

* `$DEPLOYMENT_ID` is the ID for the deployment that was created in the previous [create a deployment](/reference/ecctl-example-create-deployment.md) example.

```json
{
  "id": "20e174f6800c55261e4dfcc278b6a004",
  "name": "My second ecctl deployment",
  "resources": [
    {
      "cloud_id": "My_second_ecctl_deployment:dXMtY2VudHJhbDEuZ2NwLmZvdW5kaXQubm8kYjc0OWU2ZWExN2Y4NDg5Yzg4Y2UyOTVjZTA4ZDVjNWUkNTliZWJiYjE3ZmFkNDk2MWEwMmNkMDRmNzYyOWYxMTk=",
      "id": "b749e6ea17f8489c88ce295ce08d5c5e",
      "kind": "elasticsearch",
      "ref_id": "main-elasticsearch",
      "region": "gcp-us-central1"
    },
    {
      "elasticsearch_cluster_ref_id": "main-elasticsearch",
      "id": "59bebbb17fad4961a02cd04f7629f119",
      "kind": "kibana",
      "ref_id": "main-kibana",
      "region": "gcp-us-central1"
    },
    {
      "elasticsearch_cluster_ref_id": "main-elasticsearch",
      "id": "1ec19461253c4175a2cea6b3ccc399a8",
      "kind": "apm",
      "ref_id": "main-apm",
      "region": "gcp-us-central1"
    }
  ]
}
```

