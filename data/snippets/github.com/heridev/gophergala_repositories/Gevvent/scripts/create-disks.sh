gcloud compute disks create --size="200GB" --zone=europe-west1-c rethinkdb-disk
gcloud compute instances attach-disk --zone=europe-west1-c --disk=rethinkdb-disk --device-name temp-data kubernetes-master
gcloud compute ssh --zone=europe-west1-c kubernetes-master \
  --command "sudo rm -rf /mnt/tmp && sudo mkdir /mnt/tmp && sudo /usr/share/google/safe_format_and_mount /dev/disk/by-id/google-temp-data /mnt/tmp"
gcloud compute instances detach-disk --zone=europe-west1-c --disk rethinkdb-disk kubernetes-master

gcloud compute disks create --size="200GB" --zone=europe-west1-c consul-disk
gcloud compute instances attach-disk --zone=europe-west1-c --disk=consul-disk --device-name temp-data kubernetes-master
gcloud compute ssh --zone=europe-west1-c kubernetes-master \
  --command "sudo rm -rf /mnt/tmp && sudo mkdir /mnt/tmp && sudo /usr/share/google/safe_format_and_mount /dev/disk/by-id/google-temp-data /mnt/tmp"
gcloud compute instances detach-disk --zone=europe-west1-c --disk consul-disk kubernetes-master

