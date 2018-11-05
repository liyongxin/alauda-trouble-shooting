#!/usr/bin/env bash
set -e
set -x
make docker
docker tag alauda-trouble-shooting:master index.alauda.cn/yxli/alauda-trouble-shooting
docker push index.alauda.cn/yxli/alauda-trouble-shooting
imageId=$(docker images |awk '{print $3}'|sed -n '2p')
deleteImageId=$(docker images |awk '{print $3}'|sed -n '3p')
echo "image id is $imageId"
set +e
docker rm `docker ps -aq`
docker rmi $imageId $deleteImageId