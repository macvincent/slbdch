#update this variable with your username
username=macvincent
instance_name=go-vm3

gcloud compute scp --recurse --compress ./cache_startup.sh ./web_cache $instance_name:/home/$username
gcloud compute ssh $instance_name --command "chmod +x ./cache_startup.sh && ./cache_startup.sh" 
