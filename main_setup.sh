#!/bin/bash
username=macvincent
instance_name=main
echo "Creating instance $instance_name"
gcloud compute instances create "$instance_name" --machine-type=n1-standard-1 --image-family=debian-10 --image-project=debian-cloud --zone=us-west4-a --tags=http-server

echo "Copying files to instance $instance_name"
gcloud compute scp --recurse --compress main_startup.sh ./consistent_web_main $instance_name:/home/$username

echo "Running main_startup.sh on instance $instance_name"
gcloud compute ssh $instance_name --command "chmod +x main_startup.sh && ./main_startup.sh"

echo "SSH into instance $instance_name"
gcloud compute ssh $instance_name