#!/bin/bash
username=macvincent
instance_name=load-generator
echo "Creating instance $instance_name"
gcloud compute instances create "$instance_name" --machine-type=n1-standard-16 --image-family=debian-10 --image-project=debian-cloud --zone=us-west1-b --tags=http-server

# delay for 5 seconds to allow the instance to be created
sleep 10

echo "Copying files to instance $instance_name"
gcloud compute scp --zone=us-west1-b --recurse --compress main_startup.sh ./consistent_web_main ./load_generator $instance_name:/home/$username


echo "Running main_startup.sh on instance $instance_name"
gcloud compute ssh --zone=us-west1-b $instance_name --command "chmod +x main_startup.sh && ./main_startup.sh"

echo "SSH into instance $instance_name"
gcloud compute ssh --zone=us-west1-b $instance_name