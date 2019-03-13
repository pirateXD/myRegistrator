#!/bin/bash

ip=172.26.148.189
tag=$(cat VERSION)
id=myRegistrator
dockerRun=registrator:${tag}

docker stop $id && docker rm $id && echo "stop and remove $id"
docker run -d  --name=$id  --net=host  --volume=/var/run/docker.sock:/tmp/docker.sock  $dockerRun  \
-ip="$ip" \
-useIpFromLabel="exposeIP" \
-ttl=30 \
-ttl-refresh=10  \
etcd://qa.haidao:2379/services
