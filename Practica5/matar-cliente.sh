#!/bin/sh
sudo kubectl delete pod client
sudo kubectl apply -f pod_go.json
sleep 4
kubectl exec -ti raft-0 -- nslookup raft-0.raft.default.svc.cluster.local
kubectl exec -ti raft-0 -- nslookup raft-1.raft.default.svc.cluster.local
kubectl exec -ti raft-0 -- nslookup raft-2.raft.default.svc.cluster.local
kubectl apply -f pod_go.json
