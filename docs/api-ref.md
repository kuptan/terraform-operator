---
layout: default
title: API Reference
nav_order: 9
---

# API Reference
To view the API reference, please visit [this page](https://doc.crds.dev/github.com/kuptan/terraform-operator) 

<!-- ---
layout: default
title: API Reference
nav_order: 5
---

# API Reference

## Packages
- [run.terraform-operator.io/v1alpha1](#runterraform-operatoriov1alpha1)


## run.terraform-operator.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the run v1alpha1 API group

### Resource Types
- [Terraform](#terraform)
- [TerraformList](#terraformlist)



#### DependsOnSpec



DependsOnSpec specifies the dependency on other Terraform runs

_Appears in:_
- [TerraformSpec](#terraformspec)

| Field | Description |
| --- | --- |
| `name` _string_ | The Terraform object metadata.name |
| `namespace` _string_ | The namespace where the Terraform run exist |


#### GitSSHKeySpec





_Appears in:_
- [TerraformSpec](#terraformspec)

| Field | Description |
| --- | --- |
| `valueFrom` _[VolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#volumesource-v1-core)_ | The source of the value where the private SSH key exist |


#### Module





_Appears in:_
- [TerraformSpec](#terraformspec)

| Field | Description |
| --- | --- |
| `source` _string_ | module source, must be a valid Terraform module source |
| `version` _string_ | module version |


#### OutputSpec





_Appears in:_
- [TerraformSpec](#terraformspec)

| Field | Description |
| --- | --- |
| `key` _string_ | Output key specifies the Kubernetes secret key |
| `moduleOutputName` _string_ | The output name as defined in the source Terraform module |


#### PreviousRunStatus



PreviousRuns stores the previous run information in case the current run object was modified

_Appears in:_
- [TerraformStatus](#terraformstatus)

| Field | Description |
| --- | --- |
| `id` _string_ | Attribute name in module |


#### Terraform



Terraform is the Schema for the terraforms API

_Appears in:_
- [TerraformList](#terraformlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `run.terraform-operator.io/v1alpha1`
| `kind` _string_ | `Terraform`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[TerraformSpec](#terraformspec)_ |  |


#### TerraformList



TerraformList contains a list of Terraform



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `run.terraform-operator.io/v1alpha1`
| `kind` _string_ | `TerraformList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Terraform](#terraform) array_ |  |


#### TerraformSpec



TerraformSpec defines the desired state of Terraform

_Appears in:_
- [Terraform](#terraform)

| Field | Description |
| --- | --- |
| `terraformVersion` _string_ | The terraform version to use |
| `module` _[Module](#module)_ | The module information (source & version) |
| `backend` _string_ | A custom terraform backend configuration |
| `providersConfig` _string_ | A custom terraform providers configuration |
| `workspace` _string_ | The terraform workspae. Defaults to `default` |
| `dependsOn` _[DependsOnSpec](#dependsonspec) array_ | A list of dependencies on other Terraform runs |
| `variables` _[Variable](#variable) array_ | Variables as inputs to the Terraform module |
| `variableFiles` _[VariableFile](#variablefile) array_ | Terraform variable files |
| `outputs` _[OutputSpec](#outputspec) array_ | Terraform outputs will be written to a Kubernetes secret |
| `destroy` _boolean_ | Indicates whether a destroy job should run |
| `deleteCompletedJobs` _boolean_ | Indicates whether to keep the jobs/pods after the run is successful/completed |
| `retryLimit` _integer_ | A retry limit to be set on the Job as a backOffLimit |
| `gitSSHKey` _[GitSSHKeySpec](#gitsshkeyspec)_ | An SSH key to be able to pull modules from private git repositories |




#### Variable





_Appears in:_
- [TerraformSpec](#terraformspec)

| Field | Description |
| --- | --- |
| `key` _string_ | Terraform module variable name |
| `value` _string_ | Variable value |
| `valueFrom` _[EnvVarSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvarsource-v1-core)_ | The variable value from a key source (secret or configmap) |
| `environmentVariable` _boolean_ | EnvironmentVariable denotes if this variable should be created as environment variable |


#### VariableFile





_Appears in:_
- [TerraformSpec](#terraformspec)

| Field | Description |
| --- | --- |
| `key` _string_ | The module variable name |
| `valueFrom` _[VolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#volumesource-v1-core)_ | The source of the variable file | -->
