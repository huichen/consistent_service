package service_hash

import (
	"github.com/serialx/hashring"
	"sync"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

type ServiceHash struct {
	ringLock struct {
		sync.RWMutex
		ring *hashring.HashRing
	}
	etcdClient client.Client
}

func (hash *ServiceHash) watch(watcher client.Watcher) {
	for {
		resp, err := watcher.Next(context.Background())
		if err == nil {
			if resp.Action == "set" {
				n := resp.Node.Value
				hash.ringLock.Lock()
				hash.ringLock.ring = hash.ringLock.ring.AddNode(n)
				hash.ringLock.Unlock()
			} else if resp.Action == "delete" {
				n := resp.PrevNode.Value
				hash.ringLock.Lock()
				hash.ringLock.ring = hash.ringLock.ring.RemoveNode(n)
				hash.ringLock.Unlock()
			}
		}
	}
}

func (hash *ServiceHash) Connect(serviceName string, endPoints []string) error {
	hash.ringLock.ring = hashring.New([]string{})

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
				hash.ringLock.ring = hash.ringLock.ring.AddNode(n)
			}
		}
	}

	watcher := kapi.Watcher(serviceName, &client.WatcherOptions{Recursive: true})
	go hash.watch(watcher)
	return nil
}

func (hash *ServiceHash) Hash(key string) (string, bool) {
	hash.ringLock.RLock()
	node, ok := hash.ringLock.ring.GetNode(key)
	hash.ringLock.RUnlock()
	return node, ok
}
