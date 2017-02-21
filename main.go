package main

import (
	"log"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	checkpointInterval = time.Minute
)

func main() {
	// TODO: recover

	cp := newEndpointCheckpointer(mustNewKubeClient())

	ticker := time.NewTicker(checkpointInterval)

	for {
		select {
		case <-ticker.C:
			err := cp.checkpoint()
			if err != nil {
				log.Println("failed to checkpoint etcd endpoints:", err)
			}
			// TODO: setup iptable rules based on the endpoints.
		}
	}
}

func mustNewKubeClient() kubernetes.Interface {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	return kubernetes.NewForConfigOrDie(cfg)
}
