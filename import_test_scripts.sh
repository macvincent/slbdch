username=macvincent
instance_name=main

gcloud compute scp --recurse --compress $instance_name:/home/$username/consistent_web_main/*.txt ./experiments/kademlia_vs_cycle/kademlia/1_node
