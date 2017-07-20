package main

import (
	"fmt"
	"log"
	"strings"

	utiliptables "github.com/coreos/kenc/pkg/util/iptables"
)

const (
	natTable = "nat"

	// the services chain
	kubeServicesChain utiliptables.Chain = "KUBE-SERVICES"
	// the kubernetes postrouting chain
	kubePostroutingChain utiliptables.Chain = "KUBE-POSTROUTING"
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

// Top level chains that will not be flushed in the restore transaction.
var nonFlushChains = map[string]bool{
	"-A PREROUTING":  true,
	"-A POSTROUTING": true,
	"-A INPUT":       true,
	"-A OUTPUT":      true,
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
			ok := true
			for nc := range nonFlushChains {
				if strings.HasPrefix(line, nc) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}

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

func ensureLinkingChains(ipt utiliptables.Interface) error {
	if _, err := ipt.EnsureChain(utiliptables.TableNAT, kubeServicesChain); err != nil {
		log.Printf("Failed to ensure that %s chain %s exists: %v", utiliptables.TableNAT, kubeServicesChain, err)
		return err
	}

	tableChainsNeedJumpServices := []struct {
		table utiliptables.Table
		chain utiliptables.Chain
	}{
		{utiliptables.TableNAT, utiliptables.ChainOutput},
		{utiliptables.TableNAT, utiliptables.ChainPrerouting},
	}
	comment := "kubernetes service portals"
	args := []string{"-m", "comment", "--comment", comment, "-j", string(kubeServicesChain)}
	for _, tc := range tableChainsNeedJumpServices {
		if _, err := ipt.EnsureRule(utiliptables.Prepend, tc.table, tc.chain, args...); err != nil {
			log.Printf("Failed to ensure that %s chain %s jumps to %s: %v", tc.table, tc.chain, kubeServicesChain, err)
			return err
		}
	}

	// Create and link the kube postrouting chain.
	if _, err := ipt.EnsureChain(utiliptables.TableNAT, kubePostroutingChain); err != nil {
		log.Printf("Failed to ensure that %s chain %s exists: %v", utiliptables.TableNAT, kubePostroutingChain, err)
		return err
	}

	comment = "kubernetes postrouting rules"
	args = []string{"-m", "comment", "--comment", comment, "-j", string(kubePostroutingChain)}
	if _, err := ipt.EnsureRule(utiliptables.Prepend, utiliptables.TableNAT, utiliptables.ChainPostrouting, args...); err != nil {
		log.Printf("Failed to ensure that %s chain %s jumps to %s: %v", utiliptables.TableNAT, utiliptables.ChainPostrouting, kubePostroutingChain, err)
		return err
	}

	return nil
}
