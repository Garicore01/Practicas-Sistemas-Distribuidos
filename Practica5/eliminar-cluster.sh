#!/bin/bash
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
kubectl delete pods --all --all-namespaces
kind delete cluster