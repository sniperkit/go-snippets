#!/bin/sh
msg=`python /classify_image.pipe.py --model_dir /imagenet --image_file /images/val_100.JPEG 2> /dev/null`
uptime=`uptime`
curr_time=`date "+%Y-%m-%d %H:%M:%S.%s"`
echo "{\"system\": {\"hostname\":\"$HOSTNAME\", \"uptime\":\"$uptime\", \"curr_time\":\"$curr_time\"}, \"msg\": \"$msg\", \"arg\":$@}"
