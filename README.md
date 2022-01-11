# Terraform Operator
[![build](https://github.com/kube-champ/terraform-operator/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/kube-champ/terraform-operator/actions/workflows/build.yaml)

The Terraform Operator provides support to run Terraform modules in Kubernetes in a declaritive way as a [Kubernetes manifest](https://kubernetes.io/docs/concepts/cluster-administration/manage-deployment/).

This projects makes defining and running a Terraform module, Kubernetes native through a single Kubernetes [CRD](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/). You can run the manifest with kubectl, Terraform, GitOps tools, etc...

**Disclaimer**

This project is not a YAML to HCL converter. It just provides a way to run Terraform commands through a Kubernetes CRD. To see how this controller works underhood, have a look at the [design doc](#)

## Installation

**Helm**

```bash
  helm repo add kubechamp https://kube-champ.github.io/terraform-operator
  helm install terraform-operator kubechamp/terraform-operator
```

**Kubernetes Manifest**

```bash
  
```
- install crds
- install rbac
- install deployment 

## Usage

```yaml
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: first-module
spec:
  terraformVersion: 1.0.2

  module:
    source: IbraheemAlSaady/test/module
    version:

   ## a terraform workspace to select
  workspace:

  ## a custom terraform backend
  backend: |
    backend "local" {
      path = "/tmp/tfmodule/mytfstate.tfstate"
    }

  ## a list of terraform variables to be provided
  variables:
    - key: length
      value: "16"
    - key: AWS_ACCESS_KEY
      valueFrom:
        secretKeyRef:
          name: aws-credentials
          key: AWS_ACCESS_KEY
      environmentVariable: true

  ## files with ext '.tfvars' or '.tf' that will be mounted into the terraform runner job 
  ## to be passed to terraform as '-var-file'
  variableFiles:
    - key: terraform-env-config
      valueFrom:
        configMap:
          name: "terraform-env-config"

  ## all outputs will be written to a secret by default
  ## if provided, it will be available in the Status when you run kubectl describe run/[run-name]
  outputs:
    - key: my_new_output_name
      ## the output name from the module
      moduleOutputName: result
      ## if set, it will mask the value in the run status
      sensitive: false

  ## a flag to run a terraform destroy
  destroy: false

  ## a flag to delete the job after the job is completed
  deleteCompletedJobs: false

  ## number of retries in case of run failure
  retryLimit: 2
```