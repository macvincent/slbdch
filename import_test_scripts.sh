username=macvincent
instance_name=main
dest_dir=./experiments/kademlia_vs_cycle/cycle/
mkdir -p $dest_dir
gcloud compute scp --recurse --compress $instance_name:/home/$username/consistent_web_main/*.csv $dest_dir
