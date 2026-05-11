# Resources

Each STACKIT resource type known to `stackit-nuke` is listed below. The name in the left column is what you reference in `--include`, `--exclude`, and the `resource-types` config block.

| Name | STACKIT service | Docs |
|------|-----------------|------|
| `ComputeServer` | IaaS (compute) | [compute-server](compute-server.md) |
| `ComputeVolume` | IaaS (compute) | [compute-volume](compute-volume.md) |
| `ComputeSnapshot` | IaaS (compute) | [compute-snapshot](compute-snapshot.md) |
| `ComputeKeypair` | IaaS (compute) | [compute-keypair](compute-keypair.md) |
| `Network` | IaaS (network) | [network](network.md) |
| `Subnet` | IaaS (network) | [subnet](subnet.md) |
| `Router` | IaaS (network) | [router](router.md) |
| `SecurityGroup` | IaaS (network) | [security-group](security-group.md) |
| `FloatingIP` | IaaS (network) | [floating-ip](floating-ip.md) |
| `ObjectStorageBucket` | Object Storage | [object-storage-bucket](object-storage-bucket.md) |
| `ObjectStorageObject` | Object Storage | [object-storage-object](object-storage-object.md) |
| `SKECluster` | SKE (Kubernetes) | [ske-cluster](ske-cluster.md) |
| `PostgresFlexInstance` | PostgresFlex | [postgres-flex-instance](postgres-flex-instance.md) |
| `MongoDBFlexInstance` | MongoDBFlex | [mongodb-flex-instance](mongodb-flex-instance.md) |
| `RedisInstance` | Redis | [redis-instance](redis-instance.md) |
| `OpenSearchInstance` | OpenSearch | [opensearch-instance](opensearch-instance.md) |
| `RabbitMQInstance` | RabbitMQ | [rabbitmq-instance](rabbitmq-instance.md) |
| `LoadBalancer` | LoadBalancer | [load-balancer](load-balancer.md) |
| `DNSZone` | DNS | [dns-zone](dns-zone.md) |

Run `stackit-nuke resource-types` to print the registered names from the binary itself.
