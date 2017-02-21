package kenc

// writeNatTableRule writes a iptables NAT rule to forward the
// packets sent to the given vip to one of the given endpoints
// randomly.
// This is used to implement etcd endpoints level checkpoint.
func writeNatTableRule(vip string, endpoints []string) error {
	// TODO: implement me
	return nil
}

// saveIPtable saves iptables rule into the given file
// This is used to implement iptable level checkpoint.
func saveIPtable(filepath string) error {
	// TODO: implement me
	return nil
}

// restoreIPtableFromFile restores the iptable configuration from the give file
// that contains iptable rules.
// This is used to implement iptable level checkpoint.
func restoreIPtableFromFile(filepath string) error {
	// TODO: implement me
	return nil
}
