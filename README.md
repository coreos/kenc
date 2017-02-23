# kenc

Kuberetes etcd network checkpointer

## Checkpoint modes

- iptables

Checkpoint/restore the NAT table. Rely on kube-proxy to populate iptables rules.

```
kenc -m iptables
```

- endpoints

Checkpoint/restore etcd endpoints. Write iptables rules to iptables to ensure connectivity.

To run this mode, `kenc` MUST be started inside Kubernetes.

```
kenc -m endpoints
```
