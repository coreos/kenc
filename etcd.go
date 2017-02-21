package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/labels"
)

const (
	endpointsCheckpointFile = "endpoints.checkpoint"

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
	kubecli   kubernetes.Interface
	endpoints []string
}

func newEndpointCheckpointer(kubecli kubernetes.Interface) *endpointsCheckpointer {
	return &endpointsCheckpointer{
		kubecli: kubecli,
	}
}

func (ec *endpointsCheckpointer) checkpoint() error {
	eps, err := getEndpoints(ec.kubecli)
	if err != nil {
		return err
	}

	ec.endpoints = eps
	epss := Endpoints{
		Endpoints: eps,
	}

	b, err := json.Marshal(epss)
	if err != nil {
		return err
	}

	// TODO: return if there is no change
	f, err := os.Create(path.Join(checkpointDir, endpointsCheckpointFile))
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

// getEndpointsFromCheckpoint returns the endpoints from a previous checkpoint file.
func getEndpointsFromCheckpoint() ([]string, error) {
	b, err := ioutil.ReadFile(path.Join(checkpointDir, endpointsCheckpointFile))
	if err != nil {
		return nil, err
	}

	var eps Endpoints
	if err = json.Unmarshal(b, &eps); err != nil {
		return nil, err
	}

	return eps.Endpoints, nil
}

func getEndpoints(kubecli kubernetes.Interface) ([]string, error) {
	ls := map[string]string{
		cluterLabel: clusterName,
		appLabel:    appName,
	}

	// TODO: use client side cache
	lo := v1.ListOptions{LabelSelector: labels.SelectorFromSet(ls).String()}
	podList, err := kubecli.Core().Pods(ns).List(lo)
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
