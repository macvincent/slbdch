# Fault-Tolerant Load Balancing with Dynamic Consistent Hashing
## Abstract

Load balancing, by ensuring the equitable distribution of requests in a cluster, helps
reduce the likelihood of hotspots among cluster nodes. In this project, we will
apply a consistent hashing approach Karger et al. (1999) to perform load-balancing
tasks in a fault-tolerant manner. We will explore our ability to dynamically include
new nodes in the cluster while maintaining the consistent nature of the hash.
Performance and availability will be improved using per-shard replication similar
to that described in the ZooKeeper paper Hunt et al. (2010). We will demonstrate
the utility of our approach on a web caching system, where we will evaluate our
implementation on metrics such as latency, throughput, and the impact of nodes
going down/getting added to the system. As a stretch goal, we may also perform
evaluations on a database retrieval and update system to further demonstrate the
fault-tolerant nature of our system.


## Running Code
1. [Install Go](https://go.dev/doc/install).

1. Initialize Go modules from the project directory:
```
go mod init final_project
go mod tidy
```

2. To run server, open a new terminal and run:
```
go run .\simple_server\simple_server.go
```
3. To run client, open a new terminal and run:
```
go run .\simple_client\simple_client.go
```