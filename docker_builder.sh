#!/bin/sh
git pull && echo "pull registrator"

registryHost=$(grep "^registry" docker_run.config | awk -F'=' '{print $2}')
tag=$(grep "^version" docker_run.config | awk -F'=' '{print $2}')
if [[ -z "$tag" ]]; then
    tag="latest"
fi
echo "tag is: $tag"

docker build -t ${registryHost}registrator:${tag} -f Dockerfile .
docker push ${registryHost}dbproxy:${tag}

echo "docker Build complete!"