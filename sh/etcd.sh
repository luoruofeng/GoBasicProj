#!/bin/bash

# 启动时候更上IP参数 ./etcd 116.63.162.90

#delete docker exsiting info
c=`docker ps -q -f name="etcd"`
if [ -z "$c" ]
then
    echo "create container name:etcd"
else
    echo "delete an exsiting container:etcd and create container name:etcd"
    docker rm -f etcd
fi

n=`docker network ls | grep etcdnet`
if [ -z "$n" ]
then
    echo "create net name:etcdnet"
else
    echo "delete an exsiting net:etcdnet and create net name:etcdnet"
    docker net rm etcdnet
fi
docker network create etcdnet --driver bridge

v=`docker volume ls | grep etcd-data`
if [ -z "$v" ]
then
    echo "create volume name:etcd-data"
    docker volume create --name etcd-data
else
    echo "volume name:etcd-data is exsiting"
fi

# set env var
if [ -z "$1" ]
then
    export NODE1=127.0.0.1
else
    export NODE1=$1
fi
export DATA_DIR="etcd-data"
REGISTRY=quay.io/coreos/etcd

#create container
docker run -d --network etcdnet \
  -p 2379:2379 \
  -p 2380:2380 \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd ${REGISTRY}:latest \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://${NODE1}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${NODE1}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://${NODE1}:2380

# docker bridge测试:
# docker run -it --rm \
#     --network etcdnet \
#     --env ALLOW_NONE_AUTHENTICATION=yes \
#     bitnami/etcd:latest etcdctl --endpoints http://etcd:2379 put /message Hello
# docker run -it --rm   \
#     --network etcdnet   \
#     --env ALLOW_NONE_AUTHENTICATION=yes  \
#     bitnami/etcd:latest etcdctl --endpoints http://etcd:2379 get /message
    
# 其他网络docker测试：
# docker run -it --rm    --env ALLOW_NONE_AUTHENTICATION=yes     bitnami/etcd:latest etcdctl --endpoints http://116.63.162.90:2379 put /message2 abc
# docker run -it --rm    --env ALLOW_NONE_AUTHENTICATION=yes     bitnami/etcd:latest etcdctl --endpoints http://116.63.162.90:2379 get /message2
