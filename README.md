# Terraform Operator
[![build](https://github.com/kube-champ/terraform-operator/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/kube-champ/terraform-operator/actions/workflows/build.yaml) [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The Terraform Operator provides support to run Terraform modules in Kubernetes in a declaritive way as a [Kubernetes manifest](https://kubernetes.io/docs/concepts/cluster-administration/manage-deployment/).

This projects makes running a Terraform module, Kubernetes native through a single Kubernetes [CRD](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/). You can run the manifest with kubectl, Terraform, GitOps tools, etc...

**Disclaimer**

This project is not a YAML to HCL converter. It just provides a way to run Terraform commands through a Kubernetes CRD. To see how this controller works, have a look at the [design doc](./docs/design.md)

## Installation

**Helm**

```bash
  helm repo add kube-champ https://kube-champ.github.io/helm-charts
  helm install terraform-operator kube-champ/terraform-operator
```

Chart can be found [here](https://github.com/kube-champ/helm-charts/tree/master/charts/terraform-operator)

**Kubectl**

```bash
  kubectl apply -k https://github.com/kube-champ/terraform-operator/config/crd 
  kubectl apply -k https://github.com/kube-champ/terraform-operator/config/manifest
```

<!-- **Kubernetes Manifest**

```bash
  install crds
  install role
  install deployment
``` -->

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
    ## optional module version
    version:

   ## a terraform workspace to select
  workspace:

  ## a custom terraform backend
  backend: |
    backend "local" {
      path = "/tmp/tfmodule/mytfstate.tfstate"
    }

  ## a custom providers config
  providersConfig:

  ## a list of terraform variables to be provided
  variables:
    - key: length
      value: "16"
    
    - key: AWS_ACCESS_KEY
      valueFrom:
        ## can be configMapKeyRef as well
        secretKeyRef:
          name: aws-credentials
          key: AWS_ACCESS_KEY
      environmentVariable: true

  ## files with ext '.tfvars' or '.tf' that will be mounted into the terraform runner job 
  ## to be passed to terraform as '-var-file'
  variableFiles:
    - key: terraform-env-config
      valueFrom:
        ## can be 'secret'
        configMap:
          name: "terraform-env-config"

  ## outputs defined will be stored in a Kubernetes secret
  outputs:
    - key: my_new_output_name
      ## the output name from the module
      moduleOutputName: result

  ## a flag to run a terraform destroy
  destroy: false

  ## a flag to delete the job after the job is completed
  deleteCompletedJobs: false

  ## number of retries in case of run failure
  retryLimit: 2
```

For more examples on how to use this CRD, check the [samples](./config/samples/README.md)