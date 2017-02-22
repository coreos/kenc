package main

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	utildbus "github.com/coreos/kenc/pkg/util/dbus"
	utiliptables "github.com/coreos/kenc/pkg/util/iptables"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	utilexec "k8s.io/kubernetes/pkg/util/exec"
)

const (
	checkpointInterval = time.Minute

	vip = "10.3.0.15"

	modeEndpointsCheckpoint = "endpoints"
	modeIptablesCheckpoint  = "iptables"

	checkpointDir = "/etc/kubernetes/selfhosted-etcd"
	dirperm       = 0700
)

var (
	mode string

	// global iptables utility
	ipt utiliptables.Interface
)

func init() {
	flag.StringVar(&mode, "m", modeEndpointsCheckpoint, "kubernetes etcd netowrk checkpint mode (endpoints/iptables)")
	flag.Parse()

	ipt = utiliptables.New(utilexec.New(), utildbus.New(), utiliptables.ProtocolIpv4)
}

func main() {
	err := os.MkdirAll(checkpointDir, dirperm)
	if err != nil {
		log.Fatalf("failed to create checkpoint dir: %v", err)
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
	err := writeRouteRule(ipt, vip)
	if err != nil {
		log.Fatalf("cannot write route rule for checkpoint: %v", err)
	}

	eps, err := getEndpointsFromCheckpoint()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("cannot open endpoints checkpoint file: %v", err)
		}
	} else {
		err = writeNatTableRule(ipt, vip, eps)
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
			err = writeNatTableRule(ipt, vip, cp.endpoints)
			if err != nil {
				log.Printf("failed to update iptable rules: %v", err)
			}
		}
	}
}

func runIptablesMode() {
	err := restoreIPtablesFromFile(ipt, path.Join(checkpointDir, iptablesCheckpointFile))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to restore iptables: %v", err)
	}

	ticker := time.NewTicker(checkpointInterval)

	for {
		select {
		case <-ticker.C:
			err := saveIPtables(ipt, path.Join(checkpointDir, iptablesCheckpointFile))
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
