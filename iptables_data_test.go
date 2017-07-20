package main

var exampleTables = []byte(`# Generated by iptables-save v1.4.21 on Tue Feb 28 17:10:58 2017
*nat
:PREROUTING ACCEPT [0:0]
:INPUT ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
:POSTROUTING ACCEPT [0:0]
:DOCKER - [0:0]
:KUBE-HOSTPORTS - [0:0]
:KUBE-MARK-DROP - [0:0]
:KUBE-MARK-MASQ - [0:0]
:KUBE-NODEPORTS - [0:0]
:KUBE-POSTROUTING - [0:0]
:KUBE-SEP-5QBP3MZSYLBKSV52 - [0:0]
:KUBE-SEP-K7AB4M5TUYW23Y7Q - [0:0]
:KUBE-SEP-KFPFIE7C3VJV2765 - [0:0]
:KUBE-SEP-S3CTJS23O5GL5TRK - [0:0]
:KUBE-SEP-V226WZ5YNPSYOB5O - [0:0]
:KUBE-SEP-XDKEUXTVGH54XD6T - [0:0]
:KUBE-SERVICES - [0:0]
:KUBE-SVC-BJM46V3U5RZHCFRZ - [0:0]
:KUBE-SVC-ERIFXISQEP7F7OF4 - [0:0]
:KUBE-SVC-NPX46M4PTMTKRN6Y - [0:0]
:KUBE-SVC-TCOU7JCQXEZGVUNU - [0:0]
:KUBE-SVC-XGLOHA7QRQ3V22RZ - [0:0]
:KUBE-SVC-XP4WJ6VSLGWALMW5 - [0:0]
-A PREROUTING -m comment --comment "kube hostport portals" -m addrtype --dst-type LOCAL -j KUBE-HOSTPORTS
-A PREROUTING -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
-A OUTPUT -m comment --comment "kube hostport portals" -m addrtype --dst-type LOCAL -j KUBE-HOSTPORTS
-A OUTPUT -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
-A POSTROUTING -m comment --comment "kubernetes postrouting rules" -j KUBE-POSTROUTING
-A POSTROUTING ! -d 10.0.0.0/8 -m comment --comment "kubenet: SNAT for outbound traffic from cluster" -m addrtype ! --dst-type LOCAL -j MASQUERADE
-A POSTROUTING -s 127.0.0.0/8 -o cbr0 -m comment --comment "SNAT for localhost access to hostports" -j MASQUERADE
-A KUBE-MARK-DROP -j MARK --set-xmark 0x8000/0x8000
-A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
-A KUBE-NODEPORTS -p tcp -m comment --comment "kube-system/default-http-backend:http" -m tcp --dport 31240 -j KUBE-MARK-MASQ
-A KUBE-NODEPORTS -p tcp -m comment --comment "kube-system/default-http-backend:http" -m tcp --dport 31240 -j KUBE-SVC-XP4WJ6VSLGWALMW5
-A KUBE-POSTROUTING -m comment --comment "kubernetes service traffic requiring SNAT" -m mark --mark 0x4000/0x4000 -j MASQUERADE
-A KUBE-SEP-5QBP3MZSYLBKSV52 -s 10.216.4.2/32 -m comment --comment "kube-system/kube-dns:dns" -j KUBE-MARK-MASQ
-A KUBE-SEP-5QBP3MZSYLBKSV52 -p udp -m comment --comment "kube-system/kube-dns:dns" -m udp -j DNAT --to-destination 10.216.4.2:53
-A KUBE-SEP-K7AB4M5TUYW23Y7Q -s 10.216.1.2/32 -m comment --comment "kube-system/heapster:" -j KUBE-MARK-MASQ
-A KUBE-SEP-K7AB4M5TUYW23Y7Q -p tcp -m comment --comment "kube-system/heapster:" -m tcp -j DNAT --to-destination 10.216.1.2:8082
-A KUBE-SEP-KFPFIE7C3VJV2765 -s 10.216.6.3/32 -m comment --comment "kube-system/kubernetes-dashboard:" -j KUBE-MARK-MASQ
-A KUBE-SEP-KFPFIE7C3VJV2765 -p tcp -m comment --comment "kube-system/kubernetes-dashboard:" -m tcp -j DNAT --to-destination 10.216.6.3:9090
-A KUBE-SEP-S3CTJS23O5GL5TRK -s 104.197.104.58/32 -m comment --comment "default/kubernetes:https" -j KUBE-MARK-MASQ
-A KUBE-SEP-S3CTJS23O5GL5TRK -p tcp -m comment --comment "default/kubernetes:https" -m recent --set --name KUBE-SEP-S3CTJS23O5GL5TRK --mask 255.255.255.255 --rsource -m tcp -j DNAT --to-destination 104.197.104.58:443
-A KUBE-SEP-V226WZ5YNPSYOB5O -s 10.216.4.2/32 -m comment --comment "kube-system/kube-dns:dns-tcp" -j KUBE-MARK-MASQ
-A KUBE-SEP-V226WZ5YNPSYOB5O -p tcp -m comment --comment "kube-system/kube-dns:dns-tcp" -m tcp -j DNAT --to-destination 10.216.4.2:53
-A KUBE-SEP-XDKEUXTVGH54XD6T -s 10.216.1.13/32 -m comment --comment "kube-system/default-http-backend:http" -j KUBE-MARK-MASQ
-A KUBE-SEP-XDKEUXTVGH54XD6T -p tcp -m comment --comment "kube-system/default-http-backend:http" -m tcp -j DNAT --to-destination 10.216.1.13:8080
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.255.189/32 -p tcp -m comment --comment "kube-system/kubernetes-dashboard: cluster IP" -m tcp --dport 80 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.255.189/32 -p tcp -m comment --comment "kube-system/kubernetes-dashboard: cluster IP" -m tcp --dport 80 -j KUBE-SVC-XGLOHA7QRQ3V22RZ
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.240.10/32 -p tcp -m comment --comment "kube-system/kube-dns:dns-tcp cluster IP" -m tcp --dport 53 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.240.10/32 -p tcp -m comment --comment "kube-system/kube-dns:dns-tcp cluster IP" -m tcp --dport 53 -j KUBE-SVC-ERIFXISQEP7F7OF4
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.240.10/32 -p udp -m comment --comment "kube-system/kube-dns:dns cluster IP" -m udp --dport 53 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.240.10/32 -p udp -m comment --comment "kube-system/kube-dns:dns cluster IP" -m udp --dport 53 -j KUBE-SVC-TCOU7JCQXEZGVUNU
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.245.111/32 -p tcp -m comment --comment "kube-system/default-http-backend:http cluster IP" -m tcp --dport 80 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.245.111/32 -p tcp -m comment --comment "kube-system/default-http-backend:http cluster IP" -m tcp --dport 80 -j KUBE-SVC-XP4WJ6VSLGWALMW5
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.240.1/32 -p tcp -m comment --comment "default/kubernetes:https cluster IP" -m tcp --dport 443 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.240.1/32 -p tcp -m comment --comment "default/kubernetes:https cluster IP" -m tcp --dport 443 -j KUBE-SVC-NPX46M4PTMTKRN6Y
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.251.95/32 -p tcp -m comment --comment "kube-system/heapster: cluster IP" -m tcp --dport 80 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.251.95/32 -p tcp -m comment --comment "kube-system/heapster: cluster IP" -m tcp --dport 80 -j KUBE-SVC-BJM46V3U5RZHCFRZ
-A KUBE-SERVICES -m comment --comment "kubernetes service nodeports; NOTE: this must be the last rule in this chain" -m addrtype --dst-type LOCAL -j KUBE-NODEPORTS
-A KUBE-SVC-BJM46V3U5RZHCFRZ -m comment --comment "kube-system/heapster:" -j KUBE-SEP-K7AB4M5TUYW23Y7Q
-A KUBE-SVC-ERIFXISQEP7F7OF4 -m comment --comment "kube-system/kube-dns:dns-tcp" -j KUBE-SEP-V226WZ5YNPSYOB5O
-A KUBE-SVC-NPX46M4PTMTKRN6Y -m comment --comment "default/kubernetes:https" -m recent --rcheck --seconds 10800 --reap --name KUBE-SEP-S3CTJS23O5GL5TRK --mask 255.255.255.255 --rsource -j KUBE-SEP-S3CTJS23O5GL5TRK
-A KUBE-SVC-NPX46M4PTMTKRN6Y -m comment --comment "default/kubernetes:https" -j KUBE-SEP-S3CTJS23O5GL5TRK
-A KUBE-SVC-TCOU7JCQXEZGVUNU -m comment --comment "kube-system/kube-dns:dns" -j KUBE-SEP-5QBP3MZSYLBKSV52
-A KUBE-SVC-XGLOHA7QRQ3V22RZ -m comment --comment "kube-system/kubernetes-dashboard:" -j KUBE-SEP-KFPFIE7C3VJV2765
-A KUBE-SVC-XP4WJ6VSLGWALMW5 -m comment --comment "kube-system/default-http-backend:http" -j KUBE-SEP-XDKEUXTVGH54XD6T
COMMIT
# Completed on Tue Feb 28 17:10:58 2017
:PREROUTING ACCEPT [0:0]
# Generated by iptables-save v1.4.21 on Tue Feb 28 17:10:58 2017
*mangle
:PREROUTING ACCEPT [8551207:3181623902]
:INPUT ACCEPT [1460470:1804749150]
:FORWARD ACCEPT [7108173:1378083548]
:OUTPUT ACCEPT [1402681:247977896]
:POSTROUTING ACCEPT [8509868:1625982660]
COMMIT
# Completed on Tue Feb 28 17:10:58 2017
# Generated by iptables-save v1.4.21 on Tue Feb 28 17:10:58 2017
*filter
:INPUT DROP [0:0]
:FORWARD DROP [0:0]
:OUTPUT DROP [0:0]
:DOCKER - [0:0]
:DOCKER-ISOLATION - [0:0]
:KUBE-FIREWALL - [0:0]
:KUBE-SERVICES - [0:0]
-A INPUT -j KUBE-FIREWALL
-A INPUT -m state --state RELATED,ESTABLISHED -j ACCEPT
-A INPUT -i lo -j ACCEPT
-A INPUT -p icmp -j ACCEPT
-A INPUT -p tcp -m tcp --dport 22 -j ACCEPT
-A INPUT -p tcp -j ACCEPT
-A INPUT -p udp -j ACCEPT
-A INPUT -p icmp -j ACCEPT
-A FORWARD -j DOCKER-ISOLATION
-A FORWARD -o docker0 -j DOCKER
-A FORWARD -o docker0 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
-A FORWARD -i docker0 ! -o docker0 -j ACCEPT
-A FORWARD -i docker0 -o docker0 -j ACCEPT
-A FORWARD -p tcp -j ACCEPT
-A FORWARD -p udp -j ACCEPT
-A FORWARD -p icmp -j ACCEPT
-A OUTPUT -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
-A OUTPUT -j KUBE-FIREWALL
-A OUTPUT -m state --state NEW,RELATED,ESTABLISHED -j ACCEPT
-A OUTPUT -o lo -j ACCEPT
-A DOCKER-ISOLATION -j RETURN
-A KUBE-FIREWALL -m comment --comment "kubernetes firewall for dropping marked packets" -m mark --mark 0x8000/0x8000 -j DROP
COMMIT
# Completed on Tue Feb 28 17:10:58 2017`)

var wantTables = []byte(`*nat
:KUBE-HOSTPORTS - [0:0]
:KUBE-MARK-DROP - [0:0]
:KUBE-MARK-MASQ - [0:0]
:KUBE-NODEPORTS - [0:0]
:KUBE-POSTROUTING - [0:0]
:KUBE-SEP-5QBP3MZSYLBKSV52 - [0:0]
:KUBE-SEP-K7AB4M5TUYW23Y7Q - [0:0]
:KUBE-SEP-KFPFIE7C3VJV2765 - [0:0]
:KUBE-SEP-S3CTJS23O5GL5TRK - [0:0]
:KUBE-SEP-V226WZ5YNPSYOB5O - [0:0]
:KUBE-SEP-XDKEUXTVGH54XD6T - [0:0]
:KUBE-SERVICES - [0:0]
:KUBE-SVC-BJM46V3U5RZHCFRZ - [0:0]
:KUBE-SVC-ERIFXISQEP7F7OF4 - [0:0]
:KUBE-SVC-NPX46M4PTMTKRN6Y - [0:0]
:KUBE-SVC-TCOU7JCQXEZGVUNU - [0:0]
:KUBE-SVC-XGLOHA7QRQ3V22RZ - [0:0]
:KUBE-SVC-XP4WJ6VSLGWALMW5 - [0:0]
-A KUBE-MARK-DROP -j MARK --set-xmark 0x8000/0x8000
-A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
-A KUBE-NODEPORTS -p tcp -m comment --comment "kube-system/default-http-backend:http" -m tcp --dport 31240 -j KUBE-MARK-MASQ
-A KUBE-NODEPORTS -p tcp -m comment --comment "kube-system/default-http-backend:http" -m tcp --dport 31240 -j KUBE-SVC-XP4WJ6VSLGWALMW5
-A KUBE-POSTROUTING -m comment --comment "kubernetes service traffic requiring SNAT" -m mark --mark 0x4000/0x4000 -j MASQUERADE
-A KUBE-SEP-5QBP3MZSYLBKSV52 -s 10.216.4.2/32 -m comment --comment "kube-system/kube-dns:dns" -j KUBE-MARK-MASQ
-A KUBE-SEP-5QBP3MZSYLBKSV52 -p udp -m comment --comment "kube-system/kube-dns:dns" -m udp -j DNAT --to-destination 10.216.4.2:53
-A KUBE-SEP-K7AB4M5TUYW23Y7Q -s 10.216.1.2/32 -m comment --comment "kube-system/heapster:" -j KUBE-MARK-MASQ
-A KUBE-SEP-K7AB4M5TUYW23Y7Q -p tcp -m comment --comment "kube-system/heapster:" -m tcp -j DNAT --to-destination 10.216.1.2:8082
-A KUBE-SEP-KFPFIE7C3VJV2765 -s 10.216.6.3/32 -m comment --comment "kube-system/kubernetes-dashboard:" -j KUBE-MARK-MASQ
-A KUBE-SEP-KFPFIE7C3VJV2765 -p tcp -m comment --comment "kube-system/kubernetes-dashboard:" -m tcp -j DNAT --to-destination 10.216.6.3:9090
-A KUBE-SEP-S3CTJS23O5GL5TRK -s 104.197.104.58/32 -m comment --comment "default/kubernetes:https" -j KUBE-MARK-MASQ
-A KUBE-SEP-S3CTJS23O5GL5TRK -p tcp -m comment --comment "default/kubernetes:https" -m recent --set --name KUBE-SEP-S3CTJS23O5GL5TRK --mask 255.255.255.255 --rsource -m tcp -j DNAT --to-destination 104.197.104.58:443
-A KUBE-SEP-V226WZ5YNPSYOB5O -s 10.216.4.2/32 -m comment --comment "kube-system/kube-dns:dns-tcp" -j KUBE-MARK-MASQ
-A KUBE-SEP-V226WZ5YNPSYOB5O -p tcp -m comment --comment "kube-system/kube-dns:dns-tcp" -m tcp -j DNAT --to-destination 10.216.4.2:53
-A KUBE-SEP-XDKEUXTVGH54XD6T -s 10.216.1.13/32 -m comment --comment "kube-system/default-http-backend:http" -j KUBE-MARK-MASQ
-A KUBE-SEP-XDKEUXTVGH54XD6T -p tcp -m comment --comment "kube-system/default-http-backend:http" -m tcp -j DNAT --to-destination 10.216.1.13:8080
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.255.189/32 -p tcp -m comment --comment "kube-system/kubernetes-dashboard: cluster IP" -m tcp --dport 80 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.255.189/32 -p tcp -m comment --comment "kube-system/kubernetes-dashboard: cluster IP" -m tcp --dport 80 -j KUBE-SVC-XGLOHA7QRQ3V22RZ
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.240.10/32 -p tcp -m comment --comment "kube-system/kube-dns:dns-tcp cluster IP" -m tcp --dport 53 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.240.10/32 -p tcp -m comment --comment "kube-system/kube-dns:dns-tcp cluster IP" -m tcp --dport 53 -j KUBE-SVC-ERIFXISQEP7F7OF4
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.240.10/32 -p udp -m comment --comment "kube-system/kube-dns:dns cluster IP" -m udp --dport 53 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.240.10/32 -p udp -m comment --comment "kube-system/kube-dns:dns cluster IP" -m udp --dport 53 -j KUBE-SVC-TCOU7JCQXEZGVUNU
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.245.111/32 -p tcp -m comment --comment "kube-system/default-http-backend:http cluster IP" -m tcp --dport 80 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.245.111/32 -p tcp -m comment --comment "kube-system/default-http-backend:http cluster IP" -m tcp --dport 80 -j KUBE-SVC-XP4WJ6VSLGWALMW5
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.240.1/32 -p tcp -m comment --comment "default/kubernetes:https cluster IP" -m tcp --dport 443 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.240.1/32 -p tcp -m comment --comment "default/kubernetes:https cluster IP" -m tcp --dport 443 -j KUBE-SVC-NPX46M4PTMTKRN6Y
-A KUBE-SERVICES ! -s 10.216.0.0/14 -d 10.219.251.95/32 -p tcp -m comment --comment "kube-system/heapster: cluster IP" -m tcp --dport 80 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.219.251.95/32 -p tcp -m comment --comment "kube-system/heapster: cluster IP" -m tcp --dport 80 -j KUBE-SVC-BJM46V3U5RZHCFRZ
-A KUBE-SERVICES -m comment --comment "kubernetes service nodeports; NOTE: this must be the last rule in this chain" -m addrtype --dst-type LOCAL -j KUBE-NODEPORTS
-A KUBE-SVC-BJM46V3U5RZHCFRZ -m comment --comment "kube-system/heapster:" -j KUBE-SEP-K7AB4M5TUYW23Y7Q
-A KUBE-SVC-ERIFXISQEP7F7OF4 -m comment --comment "kube-system/kube-dns:dns-tcp" -j KUBE-SEP-V226WZ5YNPSYOB5O
-A KUBE-SVC-NPX46M4PTMTKRN6Y -m comment --comment "default/kubernetes:https" -m recent --rcheck --seconds 10800 --reap --name KUBE-SEP-S3CTJS23O5GL5TRK --mask 255.255.255.255 --rsource -j KUBE-SEP-S3CTJS23O5GL5TRK
-A KUBE-SVC-NPX46M4PTMTKRN6Y -m comment --comment "default/kubernetes:https" -j KUBE-SEP-S3CTJS23O5GL5TRK
-A KUBE-SVC-TCOU7JCQXEZGVUNU -m comment --comment "kube-system/kube-dns:dns" -j KUBE-SEP-5QBP3MZSYLBKSV52
-A KUBE-SVC-XGLOHA7QRQ3V22RZ -m comment --comment "kube-system/kubernetes-dashboard:" -j KUBE-SEP-KFPFIE7C3VJV2765
-A KUBE-SVC-XP4WJ6VSLGWALMW5 -m comment --comment "kube-system/default-http-backend:http" -j KUBE-SEP-XDKEUXTVGH54XD6T
COMMIT
`)
