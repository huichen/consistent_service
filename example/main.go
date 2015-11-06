package main

import (
	"github.com/huichen/consistent_service"
	"log"
	"time"
)

func main() {
	serviceName := "/services/busybox"
	endPoints := []string{"http://10.45.234.177:32768"}

	var service consistent_service.ConsistentService
	service.Connect(serviceName, endPoints)

	for {
		node, _ := service.GetNode("hello world")
		if node != "" {
			log.Printf("assigned to node: %s", node)
		} else {
			log.Printf("no assignment")
		}
		time.Sleep(time.Second)
	}
}
