---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-example-create-deployment.html
---

# Create a deployment [ecctl-example-create-deployment]

Let’s create a basic deployment. Elasticsearch Service supports [solutions](docs-content://get-started/index.md) and [hardware profiles](docs-content://deploy-manage/deploy/elastic-cloud/ec-configure-deployment-settings.md#ec-hardware-profiles), which pre-configure the Elastic Stack components in your deployment to best suit your particular use case. For this example, use Google Cloud Platform (GCP) to host the deployment in region `US Central 1 (Iowa)`. To know which deployment options are available by platform, see [available regions, deployment templates and instance configurations](cloud://reference/cloud-hosted/ec-regions-templates-instances.md).

Copy the following JSON payload and save it as file `create-deployment.json`.

```json
{
  "name": "My first ecctl deployment",
  "resources": {
    "elasticsearch": [
      {
        "region": "gcp-us-central1", <1>
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
              "instance_configuration_id": "gcp.data.highio.1", <2>
              "zone_count": 2, <3>
              "size": {
                "resource": "memory",
                "value": 2048 <4>
              }
            }
          ],
          "elasticsearch": {
            "version": "7.6.0" <5>
          },
          "deployment_template": {
            "id": "gcp-io-optimized" <6>
          }
        }
      }
    ],
    "kibana": [
      {
        "region": "gcp-us-central1",
        "elasticsearch_cluster_ref_id": "main-elasticsearch",
        "ref_id": "main-kibana",
        "plan": {
          "cluster_topology": [
            {
              "instance_configuration_id": "gcp.kibana.1",
              "zone_count": 1, <7>
              "size": {
                "resource": "memory",
                "value": 1024 <8>
              }
            }
          ],
          "kibana": {
            "version": "7.6.0" <9>
          }
        }
      }
    ],
    "apm": [
      {
        "region": "gcp-us-central1",
        "elasticsearch_cluster_ref_id": "main-elasticsearch",
        "ref_id": "main-apm",
        "plan": {
          "cluster_topology": [
            {
              "instance_configuration_id": "gcp.apm.1",
              "zone_count": 1, <10>
              "size": {
                "resource": "memory",
                "value": 512 <11>
              }
            }
          ],
          "apm": {
            "version": "7.6.0" <12>
          }
        }
      }
    ]
  }
}
```

1. The region for the Elasticsearch cluster
2. Instance configuration ID
3. The number of availability zones for the Elasticsearch cluster
4. The amount of memory allocated for each Elasticsearch node
5. The version of the Elasticsearch cluster
6. The template on which to base the deployment
7. The number of availability zones for Kibana
8. The amount of memory allocated for Kibana
9. The version of the Kibana instance
10. The number of availability zones for APM
11. The amount of memory allocated for APM
12. The version of the APM instance


This JSON contains the settings for a highly available Elasticsearch cluster deployed across two availability zones, a single instance of Kibana, and a single APM server.

Run the [ecctl deployment create](/reference/ecctl_deployment_create.md) command with `create-deployment.json` as a parameter. For this and other commands, you can add an optional `--track` parameter to monitor the progress.

```sh
ecctl deployment create [--track] -f create-deployment.json
```

```sh
{
  "created": true,
  "id": "7229888e7bf8350c7e4d07d7374171c0",
  "name": "My first ecctl deployment",
  "resources": [
    {
      "cloud_id": "My_first_ecctl_deployment:dXMtY2VudHJhbDEuZ2NwLmZvdW5kaXQubm8kYjFlZWVjOGQ0YWVlNGY3ZDgxNTM2Zjc1ZjZhN2Y1MDgkM2ViZTAzNmI0NDhkNDc3Y2E2ZTJjZTQ5NmE4ZDQ5ODA=",
      "credentials": {
        "password": "REDACTED",
        "username": "elastic"
      },
      "id": "b1eeec8d4aee4f7d81536f75f6a7f508",
      "kind": "elasticsearch",
      "ref_id": "main-elasticsearch",
      "region": "gcp-us-central1"
    },
    {
      "elasticsearch_cluster_ref_id": "main-elasticsearch",
      "id": "3ebe036b448d477ca6e2ce496a8d4980",
      "kind": "kibana",
      "ref_id": "main-kibana",
      "region": "gcp-us-central1"
    },
    {
      "elasticsearch_cluster_ref_id": "main-elasticsearch",
      "id": "5a03472f6dfe4f17acbe62622823b9cb",
      "kind": "apm",
      "ref_id": "main-apm",
      "region": "gcp-us-central1",
      "secret_token": "zfufcfe15eCVJk78b5"
    }
  ]
}
```

The response indicates that the request was submitted successfully. It includes the `elastic` user password, which you can use to log in to Kibana or to access the Elasticsearch REST API. Make a note of the deployment ID, which you will use in the next example.

