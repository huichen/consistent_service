package service_hash

import (
	"errors"
	"log"
	"stathat.com/c/consistent"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type ServiceHash struct {
	consis     *consistent.Consistent
	etcdClient client.Client
	connected  bool
}

func (hash *ServiceHash) watch(watcher client.Watcher) {
	for {
		resp, err := watcher.Next(context.Background())
		if err == nil {
			if resp.Action == "set" {
				n := resp.Node.Value
				hash.consis.Add(n)
			} else if resp.Action == "delete" {
				n := resp.PrevNode.Value
				hash.consis.Remove(n)
			}
		}
	}
}

func (hash *ServiceHash) Connect(serviceName string, endPoints []string) error {
	if hash.connected {
		log.Printf("Can't connected twice")
		return errors.New("math: square root of negative number")
	}

	hash.consis = consistent.New()

	cfg := client.Config{
		Endpoints:               endPoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	var err error
	hash.etcdClient, err = client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(hash.etcdClient)

	resp, err := kapi.Get(context.Background(), serviceName, nil)
	if err != nil {
		return err
	} else {
		if resp.Node.Dir {
			for _, peer := range resp.Node.Nodes {
				n := peer.Value
				hash.consis.Add(n)
			}
		}
	}

	watcher := kapi.Watcher(serviceName, &client.WatcherOptions{Recursive: true})
	go hash.watch(watcher)
	hash.connected = true
	return nil
}

func (hash *ServiceHash) Hash(key string) (string, error) {
	if !hash.connected {
		return "", errors.New("Must call connect before Hash")
	}
	node, err := hash.consis.Get(key)
	return node, err
}
