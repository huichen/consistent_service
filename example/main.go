package main

import (
	"github.com/huichen/service_hash"
	"log"
	"time"
)

func main() {
	serviceName := "/services/busybox"
	endPoints := []string{"http://10.45.234.177:32768"}

	var hash service_hash.ServiceHash
	hash.Connect(serviceName, endPoints)

	for {
		node, _ := hash.Hash("hello world")
		if node != "" {
			log.Printf("assigned to node: %s", node)
		} else {
			log.Printf("no assignment")
		}
		time.Sleep(time.Second)
	}
}
