#update this variable with your username
username=macvincent
instance_name=go-vm2

gcloud compute scp --recurse --compress ./web_cache $instance_name:/home/$username
gcloud compute ssh $instance_name --command "cd web_cache && go mod tidy && go run main.go" 
