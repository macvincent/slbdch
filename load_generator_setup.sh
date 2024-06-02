#!/bin/bash
username=macvincent
instance_name=load-generator
echo "Creating instance $instance_name"
gcloud compute instances create "$instance_name" --machine-type=n1-standard-1 --image-family=debian-10 --image-project=debian-cloud --zone=us-west4-a --tags=http-server

# delay for 5 seconds to allow the instance to be created
sleep 5

echo "Copying files to instance $instance_name"
gcloud compute scp --recurse --compress main_startup.sh ./consistent_web_main ./load_generator $instance_name:/home/$username


echo "Running main_startup.sh on instance $instance_name"
gcloud compute ssh $instance_name --command "chmod +x main_startup.sh && ./main_startup.sh"

echo "SSH into instance $instance_name"
gcloud compute ssh $instance_name