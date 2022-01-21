# Terraform Operator
[![build](https://github.com/kube-champ/terraform-operator/actions/workflows/build.yaml/badge.svg?branch=master)](https://github.com/kube-champ/terraform-operator/actions/workflows/build.yaml) [![codecov](https://codecov.io/gh/kube-champ/terraform-operator/branch/master/graph/badge.svg?token=CE594EPJOC)](https://codecov.io/gh/kube-champ/terraform-operator) [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

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
For more examples on how to use this CRD, check the [samples](./config/samples/README.md)

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
        ## can also be 'secret'
        configMap:
          name: "terraform-env-config"

  dependsOn:
    - name: run-base
      ## if its in another namespace
      namespace:
  
  ## ssh key from a secret to allow pull modules from private git repos
  gitSSHKey:
    valueFrom:
      ....

  ## outputs defined will be stored in a Kubernetes secret
  outputs:
      ## The Kubernetes Secret key
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

## Roadmap
Check the [Terraform Operator Project](https://github.com/orgs/kube-champ/projects/1) to see what's on the roadmap

## Contributing
This project welcomes contributions and suggestions. For instructions about setting up your environment to develop and extend the operator, please see [contributing.md](./docs/contributing.md)

When you submit a pull request, the pull request has to be signed