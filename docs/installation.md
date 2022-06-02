---
layout: default
title: Installation
nav_order: 2
---

# Installation
You can install the Terraform Operator either with `Helm` or directly apply the manifests with `kubectl`

**Helm**

```bash
  helm repo add kuptan https://kuptan.github.io/helm-charts
  helm install terraform-operator kuptan/terraform-operator
```

The Helm Chart source code can be found [here](https://github.com/kuptan/helm-charts/tree/master/charts/terraform-operator)

**Kubectl**

```bash
  kubectl apply -k https://github.com/kuptan/terraform-operator/config/crd 
  kubectl apply -k https://github.com/kuptan/terraform-operator/config/manifest
```
