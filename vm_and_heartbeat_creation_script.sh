#!/bin/bash
num_instances=4
username=macvincent

# Create instances
echo "Creating $num_instances instances..."
for ((i=1; i<=num_instances; i++))
do
    instance_name=go-vm$i
    gcloud compute instances create $instance_name --machine-type=n1-standard-1 --image-family=debian-10 --image-project=debian-cloud --zone=us-west4-a --tags=http-server
done

# sleep for 30 seconds
echo "Sleeping for 30 seconds..."
sleep 30

# Copy files to instances and run startup script
echo "Copying files to instances and running startup script..."
for ((i=1; i<=num_instances; i++))
do
    instance_name=go-vm$i
    gcloud compute scp --recurse --compress ./startup_script.sh ./heartbeats $instance_name:/home/$username
    nohup gcloud compute ssh $instance_name --command "chmod +x ./startup_script.sh && nohup ./startup_script.sh >> test.log" & 
done

echo "Instances created successfully!"

