#!/bin/bash

kubectl delete -f configmap.yaml
kubectl delete -f prometheus.deploy.yml
kubectl apply -f configmap.yaml
kubectl apply -f prometheus.deploy.yml