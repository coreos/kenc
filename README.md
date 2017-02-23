# kenc

Kuberetes etcd network checkpointer

## Checkpoint modes

- iptables

Checkpoint/restore the NAT table. Kenc relies on kube-proxy to populate iptables rules and ensure connectivity in this mode.

```
kenc -m iptables
```

- endpoints

Checkpoint/restore etcd endpoints. Kenc writes iptables rules to iptables to ensure connectivity periodically in this mode.

To run this mode, `kenc` MUST be started inside Kubernetes.

```
kenc -m endpoints
```
