#!/bin/sh
git pull && echo "pull registrator"

tag=$(cat VERSION)
if [ -z "$tag" ]; then
    tag="latest"
fi
echo "tag is: $tag"

docker build -t registrator:${tag} -f Dockerfile .

echo "docker Build complete!"
