// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package instances

const legacyCluster = `
{
  "associated_apm_clusters": [],
  "associated_kibana_clusters": [
    {
      "enabled": true,
      "kibana_id": "bae445a42a44e72f2f27e4b149aa496d"
    }
  ],
  "cluster_id": "b786acd298292c2d521c0e8741761b4d",
  "cluster_name": "ear",
  "elasticsearch": {
    "blocking_issues": {
      "cluster_level": [],
      "healthy": true,
      "index_level": []
    },
    "healthy": true,
    "master_info": {
      "healthy": true,
      "instances_with_no_master": [],
      "masters": [
        {
          "instances": [
            "instance-0000000024",
            "tiebreaker-0000000025",
            "instance-0000000026"
          ],
          "master_instance_name": "tiebreaker-0000000025",
          "master_node_id": "BARXKkSyRJKQUEyKeoenfw"
        }
      ]
    },
    "shard_info": {
      "available_shards": [
        {
          "instance_name": "instance-0000000024",
          "shard_count": 3
        },
        {
          "instance_name": "tiebreaker-0000000025",
          "shard_count": 3
        },
        {
          "instance_name": "instance-0000000026",
          "shard_count": 3
        }
      ],
      "healthy": true,
      "unavailable_replicas": [
        {
          "instance_name": "instance-0000000024",
          "replica_count": 0
        },
        {
          "instance_name": "tiebreaker-0000000025",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000026",
          "replica_count": 0
        }
      ],
      "unavailable_shards": [
        {
          "instance_name": "instance-0000000024",
          "shard_count": 0
        },
        {
          "instance_name": "tiebreaker-0000000025",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000026",
          "shard_count": 0
        }
      ]
    }
  },
  "healthy": true,
  "metadata": {
    "cloud_id": "ear:xxxxx",
    "endpoint": "b786acd298292c2d521c0e8741761b4d.us-east-1.aws.found.io",
    "last_modified": "2018-06-25T23:05:34.711Z",
    "version": 217
  },
  "plan_info": {
    "current": {
      "attempt_end_time": "2018-05-25T16:38:06.560Z",
      "attempt_start_time": "2018-05-25T16:36:05.517Z",
      "healthy": true,
      "plan": {
        "cluster_topology": [
          {
            "memory_per_node": 2048,
            "node_configuration": "highio.legacy",
            "node_count_per_zone": 1,
            "node_type": {
              "data": true,
              "ingest": false,
              "master": true,
              "ml": false
            }
          }
        ],
        "tiebreaker_topology": {
          "memory_per_node": 1024
        },
        "zone_count": 2
      },
      "plan_attempt_id": "96c66c73-b41c-4ba2-a155-3d9713ef4348",
      "plan_attempt_log": [],
      "plan_end_time": "0001-01-01T00:00:00.000Z"
    },
    "healthy": true,
    "history": []
  },
  "region": "us-east-1",
  "snapshots": {
    "count": 100,
    "healthy": true,
    "latest_end_time": "2018-08-10T11:07:00.551Z",
    "latest_status": "SUCCESS",
    "latest_successful": true,
    "latest_successful_end_time": "2018-08-10T11:07:00.551Z",
    "scheduled_time": "2018-08-10T11:37:00.551Z"
  },
  "status": "started",
  "system_alerts": [],
  "topology": {
    "healthy": true,
    "instances": [
      {
        "allocator_id": "i-0329cc39c8778df87",
        "container_started": true,
        "disk": {
          "disk_space_available": 49152,
          "disk_space_used": 30
        },
        "healthy": true,
        "instance_configuration": {
          "id": "aws.highio.classic",
          "name": "highio.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000024",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 2048,
          "memory_pressure": 16
        },
        "service_id": "C9FdW1aYTmKXdmYIQzXJUg",
        "service_roles": [
          "master",
          "data"
        ],
        "service_running": true,
        "service_version": "2.4.5",
        "zone": "us-east-1a"
      },
      {
        "allocator_id": "i-02f1a468dc26a6167",
        "container_started": true,
        "disk": {
          "disk_space_available": 24576,
          "disk_space_used": 0
        },
        "healthy": true,
        "instance_configuration": {
          "id": "aws.master.classic",
          "name": "master.classic",
          "resource": "memory"
        },
        "instance_name": "tiebreaker-0000000025",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 1024,
          "memory_pressure": 67
        },
        "service_id": "BARXKkSyRJKQUEyKeoenfw",
        "service_roles": [
          "master"
        ],
        "service_running": true,
        "service_version": "2.4.5",
        "zone": "us-east-1e"
      },
      {
        "allocator_id": "i-0df52ea87b711f0db",
        "container_started": true,
        "disk": {
          "disk_space_available": 49152,
          "disk_space_used": 30
        },
        "healthy": true,
        "instance_configuration": {
          "id": "aws.highio.classic",
          "name": "highio.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000026",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 2048,
          "memory_pressure": 27
        },
        "service_id": "jBuNgfu7SaSnJfA-M6qcMg",
        "service_roles": [
          "master",
          "data"
        ],
        "service_running": true,
        "service_version": "2.4.5",
        "zone": "us-east-1b"
      }
    ]
  }
}
`

const legacyECECluster = `
{
  "cluster_name": "marc-test",
  "associated_apm_clusters": [],
  "plan_info": {
    "healthy": true,
    "current": {
      "plan_attempt_id": "05473558-2df6-49ad-9469-e7ea64a67e91",
      "attempt_end_time": "2018-08-14T10:06:05.224Z",
      "source": {
        "user_id": null,
        "facilitator": "adminconsole",
        "date": "2018-08-14T10:03:58.876Z",
        "admin_id": "root",
        "action": "elasticsearch.create-cluster",
        "remote_addresses": ["192.168.255.6"]
      },
      "plan_attempt_log": [],
      "attempt_start_time": "2018-08-14T10:03:59.435Z",
      "healthy": true,
      "plan": {
        "tiebreaker_topology": {
          "memory_per_node": 1024
        },
        "elasticsearch": {
          "enabled_built_in_plugins": [],
          "user_bundles": [],
          "version": "5.6.10",
          "system_settings": {
            "enable_close_index": false,
            "use_disk_threshold": false,
            "monitoring_collection_interval": -1,
            "monitoring_history_duration": "7d",
            "destructive_requires_name": false,
            "reindex_whitelist": [],
            "auto_create_index": true,
            "watcher_trigger_engine": "scheduler",
            "scripting": {
              "inline": {
                "enabled": true,
                "sandbox_mode": true
              },
              "expressions_enabled": true,
              "stored": {
                "enabled": true,
                "sandbox_mode": true
              },
              "file": {
                "enabled": false
              },
              "mustache_enabled": true,
              "painless_enabled": true
            },
            "http": {
              "compression": true,
              "cors_enabled": false,
              "cors_max_age": 1728000,
              "cors_allow_credentials": false
            }
          },
          "user_plugins": []
        },
        "zone_count": 1,
        "cluster_topology": [{
          "node_type": {
            "master": true,
            "data": true,
            "ingest": true,
            "ml": false
          },
          "memory_per_node": 1024,
          "node_count_per_zone": 1,
          "node_configuration": "default"
        }]
      }
    },
    "history": []
  },
  "snapshots": {
    "healthy": true,
    "count": 0
  },
  "associated_kibana_clusters": [],
  "elasticsearch": {
    "healthy": true,
    "shard_info": {
      "healthy": true,
      "available_shards": [{
        "instance_name": "instance-0000000000",
        "shard_count": 0
      }],
      "unavailable_shards": [{
        "instance_name": "instance-0000000000",
        "shard_count": 0
      }],
      "unavailable_replicas": [{
        "instance_name": "instance-0000000000",
        "replica_count": 0
      }]
    },
    "master_info": {
      "healthy": true,
      "masters": [{
        "master_node_id": "51nsMZE9T_2qow-Z__w4OA",
        "master_instance_name": "instance-0000000000",
        "instances": ["instance-0000000000"]
      }],
      "instances_with_no_master": []
    },
    "blocking_issues": {
      "healthy": true,
      "cluster_level": [],
      "index_level": []
    }
  },
  "links": {

  },
  "healthy": true,
  "system_alerts": [],
  "status": "started",
  "topology": {
    "healthy": true,
    "instances": [{
      "disk": {
        "disk_space_available": 32768,
        "disk_space_used": 0
      },
      "maintenance_mode": false,
      "service_running": true,
      "healthy": true,
      "instance_name": "instance-0000000000",
      "service_version": "5.6.10",
      "service_roles": ["master", "data", "ingest"],
      "allocator_id": "192.168.44.11",
      "service_id": "51nsMZE9T_2qow-Z__w4OA",
      "zone": "ece-zone-1",
      "instance_configuration": {
        "id": "data.default",
        "name": "data.default",
        "resource": "memory"
      },
      "container_started": true,
      "memory": {
        "instance_capacity": 1024,
        "memory_pressure": 10
      }
    }]
  },
  "metadata": {
    "version": 6,
    "last_modified": "2018-08-14T10:06:05.219Z",
    "endpoint": "269d4dc158f4457491572a90ac124e6f.192.168.44.10.ip.es.io",
    "cloud_id": "marc-test:MTkyLjE2OC40NC4xMC5pcC5lcy5pbyQyNjlkNGRjMTU4ZjQ0NTc0OTE1NzJhOTBhYzEyNGU2ZiQ="
  },
  "cluster_id": "269d4dc158f4457491572a90ac124e6f"
}
`

const dntCluster = `
{
  "associated_apm_clusters": [],
  "associated_kibana_clusters": [
    {
      "enabled": true,
      "kibana_id": "1933da8d46f85380f48d9aa28cf47762"
    }
  ],
  "cluster_id": "ef97cb1bee75971e19be2522eca6a021",
  "cluster_name": "ef97cb1bee75971e19be2522eca6a021",
  "elasticsearch": {
    "blocking_issues": {
      "cluster_level": [],
      "healthy": true,
      "index_level": []
    },
    "healthy": true,
    "master_info": {
      "healthy": true,
      "instances_with_no_master": [],
      "masters": [
        {
          "instances": [
            "instance-0000000013",
            "instance-0000000009",
            "instance-0000000012",
            "instance-0000000010",
            "instance-0000000011",
            "instance-0000000008",
            "instance-0000000014"
          ],
          "master_instance_name": "instance-0000000010",
          "master_node_id": "KmuIH3JNTbuMtMSPUdF3nw"
        }
      ]
    },
    "shard_info": {
      "available_shards": [
        {
          "instance_name": "instance-0000000013",
          "shard_count": 108
        },
        {
          "instance_name": "instance-0000000009",
          "shard_count": 108
        },
        {
          "instance_name": "instance-0000000012",
          "shard_count": 108
        },
        {
          "instance_name": "instance-0000000010",
          "shard_count": 108
        },
        {
          "instance_name": "instance-0000000011",
          "shard_count": 108
        },
        {
          "instance_name": "instance-0000000008",
          "shard_count": 108
        },
        {
          "instance_name": "instance-0000000014",
          "shard_count": 108
        }
      ],
      "healthy": true,
      "unavailable_replicas": [
        {
          "instance_name": "instance-0000000013",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000009",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000012",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000010",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000011",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000008",
          "replica_count": 0
        },
        {
          "instance_name": "instance-0000000014",
          "replica_count": 0
        }
      ],
      "unavailable_shards": [
        {
          "instance_name": "instance-0000000013",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000009",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000012",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000010",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000011",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000008",
          "shard_count": 0
        },
        {
          "instance_name": "instance-0000000014",
          "shard_count": 0
        }
      ]
    }
  },
  "elasticsearch_monitoring_info": {
    "destination_cluster_ids": [
      "ef97cb1bee75971e19be2522eca6a021"
    ],
    "healthy": true,
    "last_modified": "2018-08-10T11:23:53.150Z",
    "last_update_status": "Successfully changed Marvel configuration",
    "source_cluster_ids": [
      "ef97cb1bee75971e19be2522eca6a021"
    ]
  },
  "healthy": true,
  "metadata": {
    "cloud_id": "01658f:xxxxxxxx",
    "endpoint": "ef97cb1bee75971e19be2522eca6a021.europe-west1.gcp.cloud.es.io",
    "last_modified": "2018-07-27T09:58:12.910Z",
    "version": 60
  },
  "plan_info": {
    "current": {
      "attempt_end_time": "2018-07-27T09:58:11.830Z",
      "attempt_start_time": "2018-07-27T09:47:44.708Z",
      "healthy": true,
      "plan": {
        "cluster_topology": [
          {
            "instance_configuration_id": "gcp.highio.classic",
            "memory_per_node": 16384,
            "node_count_per_zone": 1,
            "node_type": {
              "data": true,
              "ingest": true,
              "master": false,
              "ml": false
            },
            "zone_count": 2
          },
          {
            "instance_configuration_id": "gcp.ml.1",
            "node_type": {
              "data": false,
              "ingest": false,
              "master": false,
              "ml": true
            },
            "size": {
              "resource": "memory",
              "value": 32768
            },
            "zone_count": 2
          },
          {
            "instance_configuration_id": "gcp.master.classic",
            "node_type": {
              "data": false,
              "ingest": false,
              "master": true,
              "ml": false
            },
            "size": {
              "resource": "memory",
              "value": 1024
            },
            "zone_count": 3
          }
        ],
        "elasticsearch": {
          "enabled_built_in_plugins": null,
          "user_bundles": null,
          "user_plugins": null,
          "version": "6.3.2"
        },
        "tiebreaker_topology": {
          "memory_per_node": 1024
        },
        "transient": {
          "plan_configuration": {
            "calm_wait_time": 5,
            "extended_maintenance": false,
            "max_snapshot_attempts": 3,
            "move_allocators": [],
            "move_instances": [],
            "override_failsafe": false,
            "preferred_allocators": [],
            "reallocate_instances": false,
            "skip_data_migration": false,
            "skip_post_upgrade_steps": false,
            "skip_snapshot": false,
            "skip_upgrade_checker": false,
            "timeout": 131072
          },
          "strategy": {
            "grow_and_shrink": {}
          }
        }
      },
      "plan_attempt_id": "d0c2bbf8-9df1-48c4-8b49-2be49632b500",
      "plan_attempt_log": [],
      "plan_end_time": "0001-01-01T00:00:00.000Z",
      "source": {
        "action": "elasticsearch.update-cluster-plan",
        "admin_id": "proxy",
        "date": "2018-07-27T09:47:44.614Z",
        "facilitator": "adminconsole",
        "remote_addresses": [
          "213.165.163.222",
          "54.196.54.128"
        ],
        "user_id": "1372859948"
      }
    },
    "healthy": true,
    "history": []
  },
  "region": "gcp-europe-west1",
  "snapshots": {
    "count": 100,
    "healthy": true,
    "latest_end_time": "2018-08-13T00:30:32.696Z",
    "latest_status": "SUCCESS",
    "latest_successful": true,
    "latest_successful_end_time": "2018-08-13T00:30:32.696Z",
    "scheduled_time": "2018-08-13T01:00:32.696Z"
  },
  "status": "started",
  "system_alerts": [],
  "topology": {
    "healthy": true,
    "instances": [
      {
        "allocator_id": "gi-4279333926502355196",
        "container_started": true,
        "disk": {
          "disk_space_available": 393216,
          "disk_space_used": 69627
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.highio.classic",
          "name": "highio.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000013",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 16384,
          "memory_pressure": 36
        },
        "service_id": "ZL0rFFlyRmOb-1OO7hIB2Q",
        "service_roles": [
          "data",
          "ingest"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-b"
      },
      {
        "allocator_id": "gi-5212961912897534338",
        "container_started": true,
        "disk": {
          "disk_space_available": 65536,
          "disk_space_used": 0
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.ml.1",
          "name": "gcp.ml.1",
          "resource": "memory"
        },
        "instance_name": "instance-0000000009",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 32768,
          "memory_pressure": 7
        },
        "service_id": "_zL3Um1MQsyS9_eFnqlWkw",
        "service_roles": [
          "ml"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-c"
      },
      {
        "allocator_id": "gi-1485042128032659273",
        "container_started": true,
        "disk": {
          "disk_space_available": 24576,
          "disk_space_used": 0
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.master.classic",
          "name": "master.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000012",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 1024,
          "memory_pressure": 47
        },
        "service_id": "dcjWLMScQvickiZzFV5Yew",
        "service_roles": [
          "master"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-c"
      },
      {
        "allocator_id": "gi-9131968623871726719",
        "container_started": true,
        "disk": {
          "disk_space_available": 24576,
          "disk_space_used": 0
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.master.classic",
          "name": "master.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000010",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 1024,
          "memory_pressure": 65
        },
        "service_id": "KmuIH3JNTbuMtMSPUdF3nw",
        "service_roles": [
          "master"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-b"
      },
      {
        "allocator_id": "gi-4034315774680881269",
        "container_started": true,
        "disk": {
          "disk_space_available": 24576,
          "disk_space_used": 0
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.master.classic",
          "name": "master.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000011",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 1024,
          "memory_pressure": 46
        },
        "service_id": "rcly0K6bQpCoMFWtiwynzA",
        "service_roles": [
          "master"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-d"
      },
      {
        "allocator_id": "gi-4233528005825512378",
        "container_started": true,
        "disk": {
          "disk_space_available": 65536,
          "disk_space_used": 0
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.ml.1",
          "name": "gcp.ml.1",
          "resource": "memory"
        },
        "instance_name": "instance-0000000008",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 32768,
          "memory_pressure": 7
        },
        "service_id": "DGtnJIrJQxyLy-RE9VnPww",
        "service_roles": [
          "ml"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-d"
      },
      {
        "allocator_id": "gi-7231486973101299106",
        "container_started": true,
        "disk": {
          "disk_space_available": 393216,
          "disk_space_used": 69597
        },
        "healthy": true,
        "instance_configuration": {
          "id": "gcp.highio.classic",
          "name": "highio.classic",
          "resource": "memory"
        },
        "instance_name": "instance-0000000014",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 16384,
          "memory_pressure": 9
        },
        "service_id": "bcVGj-sOQTSmWMA53m0SiQ",
        "service_roles": [
          "data",
          "ingest"
        ],
        "service_running": true,
        "service_version": "6.3.2",
        "zone": "europe-west1-d"
      }
    ]
  }
}
`
