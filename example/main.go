package main

import (
	"flag"
	"github.com/huichen/consistent_service"
	"log"
	"strings"
	"time"
)

var (
	endPoints   = flag.String("endpoints", "", "Comma-separated endpoints of your etcd cluster, each starting with http://.")
	serviceName = flag.String("service_name", "", "Name of your service in etcd.")
)

func main() {
	flag.Parse()

	ep := strings.Split(*endPoints, ",")
	if len(ep) == 0 {
		log.Fatal("Can't parse --endpoints")
	}

	if *serviceName == "" {
		log.Fatal("--service_name can't be empty")
	}

	var service consistent_service.ConsistentService
	err := service.Connect(*serviceName, ep)
	if err != nil {
		log.Fatal(err)
	}

	for {
		nodes, _ := service.GetNodes("hello world", 2)
		if nodes != nil {
			log.Printf("assigned to node: %v", nodes)
		} else {
			log.Printf("no assignment")
		}
		time.Sleep(time.Second)
	}
}
