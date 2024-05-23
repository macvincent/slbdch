#!/bin/bash
num_instances=4
for ((i=1; i<=num_instances; i++))
do
    gcloud compute instances create go-vm$i --machine-type=n1-standard-1 --image-family=debian-10 --image-project=debian-cloud --zone=us-west4-a --tags=http-server
done