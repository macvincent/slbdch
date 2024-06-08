# Scalable Load Balanced Web Cache with Dynamic Consistent Hashing
## Abstract

Load balancing is crucial for preventing hotspots in a cluster by ensuring equitable distribution of requests among nodes. In this project, we develop a scalable web caching service using consistent hashing. We compare the hash chord and Kademlia tree techniques of consistent hashing and report the comparisons. Additionally, we investigate our system's capability to dynamically add/remove nodes in the cluster while maintaining the consistency of the hash. Our system also makes efforts to reduce hot URL bottlenecks by identifying hot URLs and distributing client requests for hot URLs evenly. Our evaluation focuses on metrics such as latency, throughput, and cache hit rate. We observe that our system is highly scalable and more efficient compared to traditional hashing. We also provide a qualitative analysis of the system from a distributed systems perspective. 

## Install Dependencies
1. [Install Go](https://go.dev/doc/install).


## [Optional] Running Web Caches in GCP:

1. Create the main server using main_setup.sh and leave the terminal open.
```
bash ./main_setup.sh <gcloud_username>
```

2. Next, through your google cloud console determine IP address using the steps [here](https://cloud.google.com/compute/docs/instances/view-ip-address). Make sure to explicitly set an External IP address if not present. Change `heartbeat.go`'s variable `mainAddr` to that IP address.

3. In a new terminal, run `vm_and_heartbeat_creation_script.sh` to initialize web server nodes and heart beat updates.
```
bash ./vm_and_heartbeat_creation_script.sh <number_of_instances> <gcloud_username>
```

4. Run `local_setup.sh` script to start up a web server with changed username and instance name for every cache you wish to create.
```
bash ./local_setup.sh <username> <instance_name>
```

5. Go to the `consistent_web_main/main.go` code in your main terminal and change `nodeList` to have all the IP addresses of the servers that are being used with the desired number of replicas.


# Running code on local machine
## Running consistent web cache

1. In a new termainal, change current terminal directory to `web_cache`:
```
cd web_cache
```
2. Change the variable port for every web_cache server you wish to add. Note you will need to keep track of the ports you used as you will give this information to the main server.

3. Initialize go modules from the project directory:
```
go mod tidy
```
3. Run main script:
```
go run ./
```

## Running consistent web cache heartbeat

1. In a new termainal, change current terminal directory to `heartbeats`:
```
cd heartbeats
```

3. Initialize go modules from the project directory:
2. Change the variable port for every web_cache server you wish to add. Note you will need to keep track of the ports you used as you will give this information to the main server.

3. Initialize go modules from the project directory:
```
go mod tidy
```
3. Run main script:
```
go run ./
```

## Running consistent web main (master node)
1. Update  `consistent_web_main/main.go` with a list of available ports in the object nodeList.

2. In a new termainal, change current terminal directory to `consistent_web_main`:
```
cd consistent_web_main
```
3. Initialize go modules from the project directory:
```
go mod tidy
```
4. Run main script while keeping track of open ports:
```
go run ./
```


# Running load generator
1. Run from the root directory:
```
go run load_generator/load_tester.go
```

# Dynamic node insertion/deletion
To insert a new worker node, from the master node, run:
```
go run admin/insert_remove_nodes.go insert <ip_address> <number of virtual nodes>
```
Make sure to start the web cache on the new worker and send heartbeats to master node

To remove a worker node, one way is to stop the heartbeats and the node will be removed eventually. To remove it immediately, from the master node, run:
```
go run admin/insert_remove_nodes.go remove <ip_address>
```

# Hot URLs
Configure threshold and k (gamma) in consistent_web_main/main.go
