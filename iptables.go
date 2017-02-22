package main

import (
	"io"
	"io/ioutil"
	"os"

	utiliptables "k8s.io/kubernetes/pkg/util/iptables"
)

const (
	iptablesCheckpointFile = "iptables.checkpoint"
)

// writeNatTableRule writes a iptables NAT rule to forward the
// packets sent to the given vip to one of the given endpoints
// randomly.
// This is used to implement etcd endpoints level checkpoint.
func writeNatTableRule(ipt utiliptables.Interface, vip string, endpoints []string) error {
	// TODO: implement me
	return nil
}

// saveIPtable saves iptables rule into the given file
// This is used to implement iptable level checkpoint.
func saveIPtables(ipt utiliptables.Interface, filepath string) error {
	b, err := ipt.SaveAll()
	if err != nil {
		return err
	}

	// TODO: create temp file to make save atomic
	f, err := os.Create(filepath)
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

// restoreIPtableFromFile restores the iptable configuration from the give file
// that contains iptable rules.
// This is used to implement iptable level checkpoint.
func restoreIPtablesFromFile(ipt utiliptables.Interface, filepath string) error {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// do not overwrite existing rules, do not restore counters
	return ipt.RestoreAll(b, utiliptables.NoFlushTables, utiliptables.NoRestoreCounters)
}
