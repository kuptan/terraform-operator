---
layout: default
title: Installation
nav_order: 2
---

# Installation
You can install the Terraform Operator either with `Helm` or directly apply the manifests with `kubectl`

**Helm**

```bash
  helm repo add kube-champ https://kube-champ.github.io/helm-charts
  helm install terraform-operator kube-champ/terraform-operator
```

The Helm Chart source code can be found [here](https://github.com/kube-champ/helm-charts/tree/master/charts/terraform-operator)

**Kubectl**

```bash
  kubectl apply -k https://github.com/kube-champ/terraform-operator/config/crd 
  kubectl apply -k https://github.com/kube-champ/terraform-operator/config/manifest
```
