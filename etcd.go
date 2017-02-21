package kenc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/labels"
)

const (
	endpointsCheckpointDir   = "/etc/kubernetes/selfhosted-etcd"
	dirperm                  = 0700
	endpointspCheckpointFile = "endpoints"

	ns = "kube-system"

	cluterLabel = "etcd_cluster"
	clusterName = "kube-etcd"
	appLabel    = "app"
	appName     = "etcd"
	clientPort  = "2379"
)

type Endpoints struct {
	Endpoints []string `json:"endpoints"`
}

type endpointsCheckpointer struct {
	k8s kubernetes.Interface
}

func (ec *endpointsCheckpointer) checkpoint() error {
	eps, err := getEndpoints(ec.k8s)
	if err != nil {
		return err
	}

	epss := Endpoints{
		Endpoints: eps,
	}

	b, err := json.Marshal(epss)
	if err != nil {
		return err
	}

	// TODO: return if there is no change
	err = os.MkdirAll(endpointsCheckpointDir, dirperm)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(endpointsCheckpointDir, endpointspCheckpointFile))
	if err != nil {
		return err
	}

	n, err := f.Write(b)
	if err == nil && n < len(b) {
		return io.ErrShortWrite
	}
	if err != nil {
		return err
	}

	return f.Sync()
}

func getEndpoints(k8s kubernetes.Interface) ([]string, error) {
	ls := map[string]string{
		cluterLabel: clusterName,
		appLabel:    appName,
	}

	// TODO: use client side cache
	lo := v1.ListOptions{LabelSelector: labels.SelectorFromSet(ls).String()}
	podList, err := k8s.Core().Pods(ns).List(lo)
	if err != nil {
		return nil, fmt.Errorf("failed to list running self hosted etcd pods: %v", err)
	}

	var endpoints []string
	for i := range podList.Items {
		pod := &podList.Items[i]

		switch pod.Status.Phase {
		case v1.PodRunning:
			endpoints = append(endpoints, pod.Status.PodIP+":"+clientPort)
		}
	}

	return endpoints, nil
}
