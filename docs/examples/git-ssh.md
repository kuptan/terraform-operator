---
layout: default
title: Git SSH Auth
parent: Examples
nav_order: 4
---

# Git SSH Authentication
You can specify module source from a private git repo, in order to authenticate, you must providate a private key to authenticate. Here is an example:

```yaml
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: terraform-basic
spec:
  terraformVersion: 1.0.2

  module:
    source: git::ssh://git@github.com/IbraheemAlSaady/terraform-module-test.git

  gitSSHKey:
    valueFrom:
      secret:
        secretName: git-ssh-key
        defaultMode: 0600
```

In the example above, the `spec.gitSSHKey` configures the SSH private key which will be picked up from a secret named `git-ssh-key`. The `defaultMode` is to set the permission to 600.