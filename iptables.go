package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	utiliptables "github.com/coreos/kenc/pkg/util/iptables"
)

const (
	iptablesCheckpointFile = "iptables.checkpoint"
)

var (
	selfHostedetcdChain = utiliptables.Chain("SELF-HOSTED-ETCD")
)

func writeRouteRule(ipt utiliptables.Interface, vip string) error {
	_, err := ipt.EnsureChain(utiliptables.TableNAT, selfHostedetcdChain)
	if err != nil {
		return err
	}
	// ensure the traffic to the vip jumps to our chain
	args := []string{
		"-p", "tcp",
		"--destination", vip,
		"--destination-port", clientPort,
		"-m", "tcp",
		"-m", "state",
		"--state", "NEW",
		"-j", string(selfHostedetcdChain),
	}

	_, err = ipt.EnsureRule(utiliptables.Prepend, utiliptables.TableNAT, utiliptables.ChainPrerouting, args...)
	if err != nil {
		return err
	}
	_, err = ipt.EnsureRule(utiliptables.Prepend, utiliptables.TableNAT, utiliptables.ChainOutput, args...)
	if err != nil {
		return err
	}

	return nil
}

// writeNatTableRule writes a iptables NAT rule to forward the
// packets sent to the given vip to one of the given endpoints
// randomly.
// This is used to implement etcd endpoints level checkpoint.
func writeNatTableRule(ipt utiliptables.Interface, vip string, endpoints []string) error {
	_, err := ipt.EnsureChain(utiliptables.TableNAT, selfHostedetcdChain)
	if err != nil {
		return err
	}

	n := len(endpoints)
	for i, e := range endpoints {
		args := []string{
			"-p", "tcp", // only change the new connections
			"-m", "tcp",
			"-m", "state",
			"--state", "New",
			"-m", "statistic",
			"--mode", "random",
			"--probability", fmt.Sprintf("%0.5f", 1.0/float64(i+1)),
			"-j", "DNAT",
			"--to-destination", e,
		}

		_, err = ipt.EnsureRule(utiliptables.Prepend, utiliptables.TableNAT, selfHostedetcdChain, args...)
		if err != nil {
			return err
		}
	}

	// remove remaining rules
	for i := n + 1; ; i++ {
		err = ipt.DeleteRule(utiliptables.TableNAT, selfHostedetcdChain, fmt.Sprintf("%d", i))
		if utiliptables.IsNotFoundError(err) {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// saveIPtable saves iptables rule related to etcd connectivity into the given file
// This is used to implement iptable level checkpoint.
func saveIPtables(ipt utiliptables.Interface, filepath string) error {
	b, err := ipt.SaveAll()
	if err != nil {
		return err
	}

	b, err = getKubeNATTableLines(b)
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
