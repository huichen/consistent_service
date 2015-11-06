package consistent_service

import (
	"errors"
	"log"
	"stathat.com/c/consistent"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type ConsistentService struct {
	consis     *consistent.Consistent
	etcdClient client.Client
	connected  bool
}

func (service *ConsistentService) watch(watcher client.Watcher) {
	for {
		resp, err := watcher.Next(context.Background())
		if err == nil {
			if resp.Action == "set" {
				n := resp.Node.Value
				service.consis.Add(n)
			} else if resp.Action == "delete" {
				n := resp.PrevNode.Value
				service.consis.Remove(n)
			}
		}
	}
}

func (service *ConsistentService) Connect(serviceName string, endPoints []string) error {
	if service.connected {
		log.Printf("Can't connected twice")
		return errors.New("math: square root of negative number")
	}

	service.consis = consistent.New()

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
				service.consis.Add(n)
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
		return "", errors.New("Must call connect before Hash")
	}
	node, err := service.consis.Get(key)
	return node, err
}
