package main

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	checkpointInterval = time.Minute

	vip = "todo"

	modeEndpointsCheckpoint = "endpoints"
	modeIptablesCheckpoint  = "iptables"

	checkpointDir = "/etc/kubernetes/selfhosted-etcd"
	dirperm       = 0700
)

var (
	mode string
)

func init() {
	flag.StringVar(&mode, "m", modeEndpointsCheckpoint, "kubernetes etcd netowrk checkpint mode (endpoints/iptables)")
	flag.Parse()
}

func main() {
	err := os.MkdirAll(checkpointDir, dirperm)
	if err != nil {
		log.Fatalf("failed to create checkpoint dir: %v", mode)
	}

	switch mode {
	case modeEndpointsCheckpoint:
		runEndpointsMode()
	case modeIptablesCheckpoint:
		runIptablesMode()
	default:
		log.Fatalf("unknown mode: %v", mode)
	}
}

func runEndpointsMode() {
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

func runIptablesMode() {
	err := restoreIPtableFromFile(path.Join(checkpointDir, iptablesCheckpointFile))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to restore iptables: %v", err)
	}

	ticker := time.NewTicker(checkpointInterval)

	for {
		select {
		case <-ticker.C:
			err := saveIPtable(path.Join(checkpointDir, iptablesCheckpointFile))
			if err != nil {
				log.Printf("failed to save iptables: %v", err)
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
