---
layout: default
title: Customization
nav_order: 6
---

# Customization
The Terraform Operator uses the [terraform-runner](https://github.com/kuptan/terraform-runner) as its terraform runner to execute terraform commands. If you don't want to use the default [terraform-runner](https://github.com/kuptan/terraform-runner), you can build your own.

To make the operator use your terraform runner, the Terraform Operator expects the following environment variables:

```
DOCKER_REGISTRY=docker.io
TERRAFORM_RUNNER_IMAGE=kubechamp/terraform-runner
TERRAFORM_RUNNER_IMAGE_TAG=0.0.4 ## <- this might be different
```

The above are the defaults that are passed to the operator. In helm, you can override these values by setting the following:

```yaml
terraformRunner:
  image:
    registry: docker.io
    repository: kubechamp/terraform-runner
    tag: "0.0.4"
```

## Building Your Runner
The runner of course must be a docker container at the end, the implementation in the container is up to you, however, there are few things to keep in mind.

When Terraform Operator creates Kubernetes jobs with the Terraform Runner, it sets some environment variables on the Terraform Runner container. For a technical view, have a look at this [code section](https://github.com/kuptan/terraform-operator/blob/master/api/v1alpha1/k8s_jobs.go#L16)


| Environment Variable     | Default value        | Description                                                                        |
|--------------------------|----------------------|------------------------------------------------------------------------------------|
| TERRAFORM_VERSION        | -                    | The Terraform version to install, its taken from the `spec.terraformVersion` field |
| OUTPUT_SECRET_NAME       | -                    | The Kubernetes secret to add the Terraform outputs                                 |
| TERRAFORM_WORKING_DIR    | `/tmp/tfmodule`      | The Terraform working directory                                                    |
| TERRAFORM_WORKSPACE      | `default`            | The Terraform workspace to use                                                     |
| TERRAFORM_DESTROY        | `false`              | Indicates whether to run a Terraform destroy                                       |
| TERRAFORM_VAR_FILES_PATH | `/tmp/tfvars`        | The path where var files will be mounted                                           |
| POD_NAMESPACE            | `metadata.namespace` | The Kubernetes namespace where the job is created                                  |

## Git SSH
If the the `spec.gitSSHKey` was provided to authenticate against private git repositories, the path to the ssh key will be `/root/.ssh/id_rsa`.

You need to add the ssh key `ssh-add /root/.ssh/id_rsa`