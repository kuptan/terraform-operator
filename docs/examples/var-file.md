---
layout: default
title: Variable Files
parent: Examples
nav_order: 5
---

# Terraform Variable File
You can specify variable files to your `Terraform` run either from a configmap or a secret. Here is an example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: var-file-data
data:
  data.tfvars: |-
    length = 50
---
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: terraform-var-files
spec:
  terraformVersion: 1.0.2

  module:
    source: IbraheemAlSaady/test/module
    version: 0.0.2

  variableFiles:
    - key: data-config
      valueFrom:
        configMap:
          name: var-file-data
```

In the example above, we created a configmap with a key called `data.tfvars` (extension in the key must either be `.tfvars` or `.tf`). In the `Terraform` object `spec.variableFiles`, we can configure the variable files to be used in the run.

The files in the configmap will be mounted in the Terraform job.