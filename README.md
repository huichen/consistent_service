# Consistent Service

## Prerequisite

1. Install docker on your machines

2. Launch an etcd cluster
  
  Do it your way or the easiest way via [dockerize_etcd](https://github.com/huichen/dockerize_etcd)

3. Install registrator on all machines in the cluster

  ```
  docker run -d --name=registrator --net=host --volume=/var/run/docker.sock:/tmp/docker.sock
    gliderlabs/registrator:latest etcd://<your etcd endpoint ip:port>/services
  ```
  
  Note: all services will be registered under etcd's /services keyspace.

## Run example

Go to example dir and change *serviceName* to your service name 
[(what's this)](http://gliderlabs.com/registrator/latest/user/services/) and etcd's *endPoints*. Then

    go run main.go
  
In another terminal, start and then stop a few containers under different ports like following

    docker run -it -p 8081:8081 busybox

Check how a request is hashed accordingly.
