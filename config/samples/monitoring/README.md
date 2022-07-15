# Monitoring
To test monitoring, ensure that you have a Kubernetes cluster running. You can use `kind`, `k3d`, or `minikube` for local testing.

You need to deploy the `kube-prometheus-stack` helm chart with the following

```bash
helm upgrade -i prometheus-stack prometheus-community/kube-prometheus-stack -f kube-prometheus-stack.yaml
```
