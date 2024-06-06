username=macvincent
instance_name=load-generator
dest_dir=./experiments/
mkdir -p $dest_dir
gcloud compute scp --zone=us-west1-b --recurse --compress $instance_name:/home/$username/load_generator/*.txt $dest_dir