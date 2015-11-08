# Consistent Service

This library provides

* service containerization with [docker](https://www.docker.com/)
* service discovery and automatic registration/deregistartion with [etcd](https://github.com/coreos/etcd) and [registrator](https://github.com/gliderlabs/registrator)
* service node assignment with [consistent hashing](https://godoc.org/stathat.com/c/consistent)

## Prerequisite

1. Install docker on your machines

2. Launch an etcd cluster
  
  Do it your way or the easiest way via [etcd_docker](https://github.com/huichen/etcd_docker)

3. Install registrator on all machines in the cluster

  ```
  docker run -d --name=registrator --net=host --volume=/var/run/docker.sock:/tmp/docker.sock
    gliderlabs/registrator etcd://<your etcd endpoint ip:port>/services
  ```
  
  Note: all services will be registered under etcd's /services keyspace.

## Run example

Start a few busybox containers with different ports

    docker run -it -p 8081:8081 busybox
    docker run -it -p 8082:8082 busybox

Go to example dir and

    go run main.go --endpoints=http://<your etcd endpoint ip:port> --service_name=/services/busybox
  
Then stop/restart a few containers and check how a service node is assigned accordingly.
