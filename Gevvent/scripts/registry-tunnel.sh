#! /bin/bash
gcloud compute ssh kubernetes-minion-1 --ssh-flag="-L $1:$2:$1" --ssh-flag="-N" --ssh-flag="-n"
gcloud compute ssh kubernetes-minion-1 --ssh-flag="-L 3001:10.0.42.251:80" --ssh-flag="-N" --ssh-flag="-n"
gcloud compute ssh kubernetes-minion-1 --ssh-flag="-L 3001:10.0.56.216:8181" --ssh-flag="-N" --ssh-flag="-n"
