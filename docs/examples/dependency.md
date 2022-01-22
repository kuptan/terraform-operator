---
layout: default
title: Terraform Dependencies
parent: Examples
nav_order: 3
---

# Dependencies
A `Terraform` object can depend on another `Terraform` object. Lets take the following example:

```yaml
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: terraform-run1
spec:
  terraformVersion: 1.0.2

  module:
    source: IbraheemAlSaady/test/module
    version: 0.0.2

  variables:
    - key: length
      value: "16"

  outputs:
    - key: result
      moduleOutputName: result
---
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: terraform-run2
spec:
  terraformVersion: 1.0.2

  module:
    source: IbraheemAlSaady/test/module
    version: 0.0.2

  dependsOn:
    - name: terraform-run1

  variables:
    - key: length
      value: "16"

  outputs:
    - key: result
      moduleOutputName: result
```

In the example above, the run `terraform-run2` will not run and will be in a waiting state until `terraform-run1` is in a completed state. 

You can also create a dependency on a Terraform object in a different namespace as follows:

```yaml
...
spec:
  ...
  dependsOn:
    - name: terraform-run1
      namespace: some_other_namespacce
```