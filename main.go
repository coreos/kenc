package main

import (
	"flag"
	"log"
	"os"
	"path"
	"time"

	utildbus "github.com/coreos/kenc/pkg/util/dbus"
	utilexec "github.com/coreos/kenc/pkg/util/exec"
	utiliptables "github.com/coreos/kenc/pkg/util/iptables"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	modeEndpointsCheckpoint = "endpoints"
	modeIptablesCheckpoint  = "iptables"

	dirperm = 0700

	defaultVIP            = "10.3.0.15"
	defaultCheckpointDir  = "/etc/kubernetes/selfhosted-etcd"
	defaultClusterInteval = 30 * time.Second
)

var (
	mode               string
	r                  bool
	vip                string
	checkpointDir      string
	checkpointInterval time.Duration

	// global iptables utility
	ipt utiliptables.Interface
)

func init() {
	flag.StringVar(&mode, "m", modeEndpointsCheckpoint, "kubernetes etcd netowrk checkpint mode (endpoints/iptables)")
	flag.BoolVar(&r, "r", false, "network recovery only")
	flag.StringVar(&vip, "etcd-service-ip", defaultVIP, "the kuberentes service ip of the etcd cluster")
	flag.StringVar(&checkpointDir, "checkpoint-dir", defaultCheckpointDir, "the directory to store/restore checkpoints")
	flag.DurationVar(&checkpointInterval, "checkpoint-interval", defaultClusterInteval, "the time interval to take checkpoints")
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
	if r {
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
		os.Exit(0)
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
	if r {
		err := restoreIPtablesFromFile(ipt, path.Join(checkpointDir, iptablesCheckpointFile))
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to restore iptables: %v", err)
		}
		os.Exit(0)
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
