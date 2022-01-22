---
layout: default
title: Terraform AWS
parent: Examples
nav_order: 1
---

# Using AWS Resources
Lets create a Terraform object that will run a module which creates an S3 bucket, module source can be found [here](https://github.com/IbraheemAlSaady/terraform-module-test/blob/main/modules/aws/main.tf)

```yaml
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: terraform-aws-s3
spec:
  terraformVersion: 1.0.2

  module:
    source: IbraheemAlSaady/test/module//modules/aws
    version: 0.0.2

  variables:
    - key: name
      value: "mys3bucket"
    - key: AWS_DEFAULT_REGION
      value: eu-west-1
      environmentVariable: true
    - key: AWS_ACCESS_KEY_ID
      environmentVariable: true
      valueFrom:
        secretKeyRef:
          name: aws-credentials
          key: AWS_ACCESS_KEY_ID
    - key: AWS_SECRET_ACCESS_KEY
      environmentVariable: true
      valueFrom:
        secretKeyRef:
          name: aws-credentials
          key: AWS_SECRET_ACCESS_KEY

  backend: |
    backend "s3" {
      bucket = "mybucket"
      key    = "path/to/my/key"
      region = "eu-west-1"
    }
  
  providersConfig: |
    terraform {
      required_providers {
        aws = {
          source  = "hashicorp/aws"
          version = "~> 3.0"
        }
      }
    }

    provider "aws" {
      region = "eu-west-1"
    }

  outputs:
    - key: bucket_id
      moduleOutputName: id 
```

As you notice, we're passing `AWS_ACCESS_KEY_ID` and  `AWS_SECRET_ACCESS_KEY` variables as `environmentVariable`. The values are picked up from a secret called `aws-credentials` which is created in the same namespace where the `Terraform` object is created. This is to authenticate the terraform AWS provider

We also provided the `providersConfig` section which configures the Terraform providers. A `backend` section is also configured.

Finally, there is only one output defined, which is `bucket_id`. A secret will be created for the run where the secret key will be `bucket_id` and the value is picked up from the module output, which is `id` as defined in the module source code.