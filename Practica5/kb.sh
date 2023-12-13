#!/bin/bash
echo -e "\nEliminando el cluster"
./eliminar-cluster.sh  #&>/dev/null
echo -e "\nCreando el cluster"
./kind-with-registry.sh #&>/dev/null

echo -e "\nEliminando ficheros residuales"
rm Dockerfiles/servidor/srvraft &>/dev/null
rm Dockerfiles/cliente/cltraft &>/dev/null

echo -e "\nCompilando los ficheros Golang"
cd raft/cmd/srvraft
CGO_ENABLED=0 go build -o ../../../Dockerfiles/servidor/srvraft main.go
cd ../../pkg/cltraft
CGO_ENABLED=0 go build -o ../../../Dockerfiles/cliente/cltraft cltraft.go

echo -e "\nCreando las imágenes Docker"
cd ../../../Dockerfiles/servidor
docker build . -t localhost:5001/servidor:latest
docker push localhost:5001/servidor:latest
cd ../cliente
docker build . -t localhost:5001/cliente:latest
docker push localhost:5001/cliente:latest

cd ../..

echo -e "\nLanzando Kubernetes"

kubectl apply -f service_go.json
kubectl apply -f statefulset_go.json
sleep 4
kubectl exec -ti raft-0 -- nslookup raft-0.raft.default.svc.cluster.local
kubectl exec -ti raft-0 -- nslookup raft-1.raft.default.svc.cluster.local
kubectl exec -ti raft-0 -- nslookup raft-2.raft.default.svc.cluster.local
kubectl apply -f pod_go.json
