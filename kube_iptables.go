package main

import (
	"fmt"
	"strings"

	utiliptables "github.com/coreos/kenc/pkg/util/iptables"
)

const (
	natTable = "nat"
)

var kubeKeywords = map[string]bool{
	// Chains defined in kube-proxy as global consts
	"KUBE-SERVICES":    true,
	"KUBE-HOSTPORTS":   true,
	"KUBE-NODEPORTS":   true,
	"KUBE-POSTROUTING": true,
	"KUBE-MARK-MASQ":   true,
	"KUBE-MARK-DROP":   true,
	// Chains/Rules defined in kube-proxy as in line consts
	// https://github.com/kubernetes/kubernetes/blob/20ed2a2744cdb0f790df9f792cdda5727726e102/pkg/proxy/iptables/proxier.go#L1315
	"KUBE-SVC-": true,
	"KUBE-SEP-": true,
	"KUBE-FW-":  true,
	"KUBE-XLB-": true,
}

func getKubeNATTableLines(save []byte) ([]byte, error) {
	var (
		lines []string
		ri    int
	)

	natTableStarts := "*" + natTable

	// find beginning of nat table and save it
	for ri < len(save) {
		line, n := utiliptables.ReadLine(ri, save)
		ri = n
		if strings.HasPrefix(line, natTableStarts) {
			lines = append(lines, line)
			break
		}
	}

	if ri >= len(save) {
		// nothing to checkpoint
		return nil, nil
	}

	var done bool
	// parse table lines
	for ri < len(save) && !done {
		line, n := utiliptables.ReadLine(ri, save)
		ri = n

		switch {
		case strings.HasPrefix(line, "COMMIT"):
			// save commit line and we are done!
			lines = append(lines, line)
			done = true
		case strings.HasPrefix(line, "*"):
			// unexpected new table before we commit the NAT table
			return nil, fmt.Errorf("unexpected table line: %v", line)
		case strings.HasPrefix(line, "#"):
			// ignore comment lines
		default:
			// normal lines, save them if they match Kube rules
			for k := range kubeKeywords {
				if strings.Contains(line, k) {
					lines = append(lines, line)
					break
				}
			}
		}
	}

	if !done {
		return nil, fmt.Errorf("failed to find the COMMIT LINE")
	}

	after := []byte(strings.Join(lines, "\n"))
	after = append(after, '\n')

	return after, nil
}
