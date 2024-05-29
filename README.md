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
1. Create a Google Cloud VM (`gcloud projects list` to view lists of projects):
```
 gcloud compute instances create go-vm1 --machine-type=n1-standard-1 --image-family=debian-10 --image-project=debian-cloud --zone=us-west4-a --tags=http-server
```
You will also need to update your GCP rules to allow external connection requests using these [instructions](https://www.geeksforgeeks.org/how-to-open-port-in-gcp-vm/).

2. Move the repo to a newly create VM:
```
gcloud compute scp --recurse go-vm1:./ ../
```
3. Connect to the VM usign command (`sudo gcloud compute config-ssh` configures ssh if first time):
```
gcloud compute ssh go-vm1 --zone=us-west4-a
```

4. In the ssh shell, make `setup.sh` an executable file:
```
chmod +x setup.sh
```

5. Run setup script:
```
./setup.sh
```

## Running consistent web main

1. In a new termainal, change current terminal directory to `web_cache`:
```
cd web_cache
```
2. Initialize go modules from the project directory:
```
go mod tidy
```
3. Run main script while keeping track of open ports:
```
go run ./
```

## Running consistent web main
1. Update  `consistent_web_main/main.go` with a list of available ports.
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
