#update this variable with your username
username=macvincent
instance_name=go-vm2

gcloud compute scp --recurse --compress startup_script.sh ./web_cache $instance_name:/home/$username
gcloud compute ssh $instance_name --command "chmod +x startup_script.sh && ./startup_script.sh" 

