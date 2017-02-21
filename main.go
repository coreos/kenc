package main

import (
	"log"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	checkpointInterval = time.Minute

	vip = "todo"
)

func main() {
	eps, err := getEndpointsFromCheckpoint()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("cannot open endpoints checkpoint file: %v", err)
		}
	} else {
		err = writeNatTableRule(vip, eps)
		if err != nil {
			log.Fatalf("cannot setup iptable rules for recovery: %v", err)
		}
	}

	cp := newEndpointCheckpointer(mustNewKubeClient())

	ticker := time.NewTicker(checkpointInterval)

	for {
		select {
		case <-ticker.C:
			err := cp.checkpoint()
			if err != nil {
				log.Printf("failed to checkpoint etcd endpoints: %v", err)
			}
			err = writeNatTableRule(vip, cp.endpoints)
			if err != nil {
				log.Printf("failed to update iptable rules: %v", err)
			}
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
