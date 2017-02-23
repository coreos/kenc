# kenc

Kuberetes etcd network checkpointer

## Checkpoint modes

- iptables

Checkpoint/restore the NAT table. Rely on kube-proxy to populate iptables rules.

- endpoints

Checkpoint/restore etcd endpoints. Write iptables rules to iptables to ensure connectivity.
