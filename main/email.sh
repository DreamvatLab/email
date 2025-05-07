#!/bin/sh
docker rm -f email
docker rmi email
docker load -i ./email.tar
docker run --name email -d --restart always --network host -v /data/email/configs.json:/app/configs.json lukiya/email