package consistent_service

import (
	"errors"
	"github.com/huichen/consistent_hashing"
	"log"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type ConsistentService struct {
	consis     *consistent_hashing.Consistent
	etcdClient client.Client
	connected  bool
	nodes      map[string]bool
}

func (service *ConsistentService) watch(watcher client.Watcher) {
	for {
		resp, err := watcher.Next(context.Background())
		if err == nil {
			if resp.Action == "set" {
				n := resp.Node.Value
				if _, ok := service.nodes[n]; !ok {
					service.consis.Add(n)
					service.nodes[n] = true
				}
			} else if resp.Action == "delete" {
				n := resp.PrevNode.Value
				if _, ok := service.nodes[n]; ok {
					service.consis.Remove(n)
					delete(service.nodes, n)
				}
			}
		}
	}
}

// serviceName is like "/services/busybox"
// endPoints is an array of "http://<etcd client ip:port>"
func (service *ConsistentService) Connect(serviceName string, endPoints []string) error {
	if service.connected {
		log.Printf("Can't connected twice")
		return errors.New("math: square root of negative number")
	}

	service.nodes = make(map[string]bool)

	service.consis = consistent_hashing.New()

	cfg := client.Config{
		Endpoints:               endPoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	var err error
	service.etcdClient, err = client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(service.etcdClient)

	resp, err := kapi.Get(context.Background(), serviceName, nil)
	if err != nil {
		return err
	} else {
		if resp.Node.Dir {
			for _, peer := range resp.Node.Nodes {
				n := peer.Value
				if _, ok := service.nodes[n]; !ok {
					service.consis.Add(n)
					service.nodes[n] = true
				}
			}
		}
	}

	watcher := kapi.Watcher(serviceName, &client.WatcherOptions{Recursive: true})
	go service.watch(watcher)
	service.connected = true
	return nil
}

func (service *ConsistentService) GetNode(key string) (string, error) {
	if !service.connected {
		return "", errors.New("Must call connect first")
	}
	node, err := service.consis.Get(key)
	return node, err
}

// Gets the N closest distinct nodes.
func (service *ConsistentService) GetNodes(key string, n int) ([]string, error) {
	if !service.connected {
		return nil, errors.New("Must call connect first")
	}
	nodes, err := service.consis.GetN(key, n)
	return nodes, err
}
