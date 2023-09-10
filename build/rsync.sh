#!/bin/bash -e
DIR=$(readlink -f "$0") && DIR=$(dirname "$DIR") && cd "$DIR"

source ./common.sh

rsync -avzP -e "ssh -i /Users/funplus/work/teaapp_key/aws_newkey.pem" ${EXE}  centos@ec2-54-251-33-204.ap-southeast-1.compute.amazonaws.com:/mnt/htdocs/cron/notification/android