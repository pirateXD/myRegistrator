#!/bin/bash

ip=172.26.147.228
tag=$(grep "^version" docker_run.config | awk -F'=' '{print $2}')
id=myRegistrator
dockerRun=qa.haidao:5000/registrator:${tag}

docker stop $id && docker rm $id && echo "stop and remove $id"
docker run -d  --name=$id  --net=host  --volume=/var/run/docker.sock:/tmp/docker.sock  $dockerRun  \
-ip="$ip" \
-useIpFromLabel="exposeIP" \
-ttl=30 \
-ttl-refresh=10  \
etcd://qa.haidao:2379/services