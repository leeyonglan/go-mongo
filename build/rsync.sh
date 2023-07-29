#!/bin/bash -e
DIR=$(readlink -f "$0") && DIR=$(dirname "$DIR") && cd "$DIR"

source ./common.sh

rsync -avz -e "ssh -i /Volumes/data/teaapp_aws_server_key/aws_newkey.pem" ${EXE}  centos@ec2-54-251-33-204.ap-southeast-1.compute.amazonaws.com:/mnt/htdocs/cron/notification/ios