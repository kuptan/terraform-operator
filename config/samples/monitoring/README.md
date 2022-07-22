# Monitoring
This is used for internal tests to scrape metrics from the controller. You will need to have the followng

* Local/Remote Kubernetes cluster: For local testing you can use Kind, K3d, or Minikube

## Deploy Prometheus Stack
Install the prometheus stack with the following

```bash
make monitoring
```