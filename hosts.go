package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

const etcdHostsFilename = "etcd-hosts.checkpoint"

type hostInfo struct {
	HostName string
	IP       string
}

func runHostsCheckpointer(kubecli kubernetes.Interface) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			hosts, err := getHosts(kubecli)
			if err != nil {
				log.Printf("failed to checkpoint etcd hosts: %v", err)
				continue
			}
			if len(hosts) == 0 {
				continue
			}
			fp := filepath.Join(checkpointDir, etcdHostsFilename)
			err = saveHostsCheckpoint(fp, hosts)
			if err != nil {
				log.Printf("failed to update etcd hosts file (%s): %v", fp, err)
			}
		}
	}
}

func getHosts(kubecli kubernetes.Interface) ([]*hostInfo, error) {
	ls := map[string]string{
		cluterLabel: clusterName,
		appLabel:    appName,
	}

	// TODO: use client side cache
	lo := metav1.ListOptions{LabelSelector: labels.SelectorFromSet(ls).String()}
	podList, err := kubecli.Core().Pods(api.NamespaceSystem).List(lo)
	if err != nil {
		return nil, fmt.Errorf("failed to list running self hosted etcd pods: %v", err)
	}

	var hs []*hostInfo
	for i := range podList.Items {
		pod := &podList.Items[i]

		switch pod.Status.Phase {
		case v1.PodRunning:
			h := &hostInfo{
				HostName: pod.Name,
				IP:       pod.Status.PodIP,
			}
			hs = append(hs, h)
		}
	}
	return hs, nil
}

func saveHostsCheckpoint(filepath string, hosts []*hostInfo) error {
	f, err := ioutil.TempFile(checkpointDir, "tmp-etcd-hosts")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	b := getHostsBytes(hosts)
	if _, err := f.Write(b); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filepath)
}

func getHostsBytes(hosts []*hostInfo) []byte {
	var buf bytes.Buffer
	for _, h := range hosts {
		buf.WriteString(fmt.Sprintf("%s %s.%s.%s.svc.cluster.local\n", h.IP, h.HostName, clusterName, api.NamespaceSystem))
	}
	return buf.Bytes()
}
