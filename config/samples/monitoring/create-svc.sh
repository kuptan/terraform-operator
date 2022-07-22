#!/bin/bash

kubectl run curl --image=alpine/curl:3.14 -- nslookup host.k3d.internal
	
sleep 15

IP_ADDRESS=$(kubectl logs curl | grep "Address: " | sed 's/.*: //g;s/ .*//g')

if [ -z "${IP_ADDRESS}" ]
then
  echo "IP address was not found, getting pod logs"
  kubectl logs curl
  kubectl delete pod/curl
  exit 1
fi

echo ">> found host IP address: ${IP_ADDRESS}"

kubectl delete pod/curl

echo ">> creating terraform operator monitoring service"

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: terraform-operator-metrics
  labels:
    app: terraform-operator
spec:
  ports:
    - name: metrics
      port: 8080
      targetPort: 8080
---
apiVersion: v1
kind: Endpoints
metadata:
  name: terraform-operator-metrics
  labels:
    app: terraform-operator
subsets:
- addresses:
  - ip: "${IP_ADDRESS}"
  ports: 
  - name: metrics
    port: 8080
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: terraform-operator-metrics
  labels:
    app: terraform-operator
spec:
  endpoints:
  - interval: 30s
    port: metrics
  namespaceSelector:
    matchNames:
    - default
  selector:
    matchLabels:
      app: terraform-operator
EOF
