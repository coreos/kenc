package kenc

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/labels"
)

const (
	ns = "kube-system"

	cluterLabel = "etcd_cluster"
	clusterName = "kube-etcd"
	appLabel    = "app"
	appName     = "etcd"
	clientPort  = "2379"
)

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
