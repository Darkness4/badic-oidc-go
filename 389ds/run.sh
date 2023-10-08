#!/bin/sh

# Create volume for data
docker volume create 389ds_data

docker run --rm -it \
  --env-file "$(pwd)/.env" \
  -p 3389:3389 \
  -v "389ds_data:/data" \
  --name 389ds \
  docker.io/389ds/dirsrv:2.4
