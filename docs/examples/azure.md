---
layout: default
title: Terraform Azure
parent: Examples
nav_order: 2
---

# Using Azure Resources
Lets create a Terraform object that will run a module which creates an S3 bucket, module source can be found [here](https://github.com/IbraheemAlSaady/terraform-module-test/blob/main/modules/azure/main.tf)

```yaml
apiVersion: run.terraform-operator.io/v1alpha1
kind: Terraform
metadata:
  name: terraform-az-storage
spec:
  terraformVersion: 1.0.2

  module:
    source: IbraheemAlSaady/test/module//modules/azure
    version: 0.0.2

  variables:
    - key: name
      value: "mystorage"
    - key: ARM_CLIENT_ID
      environmentVariable: true
      valueFrom:
        secretKeyRef:
          name: azure-credentials
          key: ARM_CLIENT_ID
    - key: ARM_CLIENT_SECRET
      environmentVariable: true
      valueFrom:
        secretKeyRef:
          name: azure-credentials
          key: ARM_CLIENT_SECRET
    - key: ARM_SUBSCRIPTION_ID
      environmentVariable: true
      valueFrom:
        secretKeyRef:
          name: azure-credentials
          key: ARM_SUBSCRIPTION_ID
    - key: ARM_TENANT_ID
      environmentVariable: true
      valueFrom:
        secretKeyRef:
          name: azure-credentials
          key: ARM_TENANT_ID

  backend: |
    backend "azurerm" {
      resource_group_name  = "StorageAccount-ResourceGroup"
      storage_account_name = "abcd1234"
      container_name       = "tfstate"
      key                  = "prod.terraform.tfstate"
    }
  
  providersConfig: |
    terraform {
      required_providers {
        azurerm = {
          source  = "hashicorp/azurerm"
          version = "=2.46.0"
        }
      }
    }

    provider "azurerm" {
      features {}
    }
  
  outputs:
    - key: storage_account_id
      moduleOutputName: id 

```

As you notice, we're passing `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_SUBSCRIPTION_ID`, and `ARM_TENANT_ID` variables as `environmentVariable`. The values are picked up from a secret called `ARM_TENANT_ID` which is created in the same namespace where the `Terraform` object is created. This is to authenticate the terraform Azure provider

We also provided the `providersConfig` section which configures the Terraform providers. A `backend` section is also configured.

Finally, there is only one output defined, which is `storage_account_id`. A secret will be created for the run where the secret key will be `storage_account_id` and the value is picked up from the module output, which is `id` as defined in the module source code.