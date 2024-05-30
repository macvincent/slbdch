# Scalable Load Balanced Web Cache with Dynamic Consistent Hashing
## Abstract

Load balancing, by ensuring the equitable distribution of requests in a cluster, helps
reduce the likelihood of hotspots among cluster nodes. In this project, we will
apply a consistent hashing approach to develop a scalable web caching service. 
We will explore our ability to dynamically include new nodes in the cluster while 
maintaining the consistent nature of the hash. We will also explore both the disk and the
Kademlia tree versions of consistent hashing. We will demonstrate the utility of our approach, 
where we will evaluate our implementation on metrics such as latency, throughput, and the 
impact of nodesgoing down/getting added to the system.

## Install Dependencies
1. [Install Go](https://go.dev/doc/install).


## [Optional] Running Web Caches in GCP:

1. Create one server for main server using main_setup.sh

2. Through your google cloud portal determine IP address using the steps [here](https://cloud.google.com/compute/docs/instances/view-ip-address). Make sure to explicitly set an External IP address if not present. Change heartbeat.go code's variable mainAddr to that IP address.

3. Run vm_creation_script.sh with changed username and instance name going from go-vm[k] where k is an integer with the following range [1, number of servers desired].

4. Run local_setup.sh script to start up web server with changed username and instance name for every server you wish to create.

5. Go to consistent_web_main/main.go code and change nodeList to have all the IP addresses of the servers that are being used with the desired number of replicas.


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
```
go mod tidy
```
3. Run main script:
```
go run ./
```

## Running consistent web main
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
