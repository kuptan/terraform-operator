# Requirements 
If you're interested in contributing to this project, you'll need:

* Go installed - see this [Getting Started](https://golang.org/doc/install) guide for Go.
* Docker installed - see this [Getting Started](https://docs.docker.com/install/) guide for Docker.
* `Kubebuilder` -  see this [Quick Start](https://book.kubebuilder.io/quick-start.html) guide for installation instructions.
* Kubernetes command-line tool `kubectl` 
* Access to a Kubernetes cluster. Some options are:
	* k3d
  * minikube
  * cloud managed (AWS EKS, Azure AKS, GKE, etc...)

# Quick start
This project uses the [terraform-runner](https://github.com/kube-champ/terraform-runner) project to execute terraform commands. 
If you're not familiar with how this controller works under the hood, its highly recommended to visit the [design docs](./design.md) first.

#### Now lets start

1. Clone the repositorry and open the project in a code editor (e.g: visual studio code)
2. In the root directory, create a `.env` file with the following environment variables
```bash
export DOCKER_REGISTRY=docker.io
export TERRAFORM_RUNNER_IMAGE=kubechamp/terraform-runner

## For the latest tags, check docker hub: https://hub.docker.com/r/kubechamp/terraform-runner
export TERRAFORM_RUNNER_IMAGE_TAG=0.0.3 # <-- might be a higher version
export KNOWN_HOSTS_CONFIGMAP_NAME=terraform-operator-known-hosts
```
3. Source the .env with `source .env`
4. Once you have a Kubernetes cluster running, create a `kubeconfig` file in the root of the project with the config of the Kubernetes cluster
5. Create the following Kubernetes RBAC objects, this is needed by the `terraform-runner` due to writing outputs to a Kubernetes secret
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: terraform-runner
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  # "namespace" omitted since ClusterRoles are not namespaced
  name: terraform-runner
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: terraform-runner
subjects:
- kind: ServiceAccount
  name: terraform-runner # name of your service account
  namespace: default # this is the namespace your service account is in
roleRef: # referring to your ClusterRole
  kind: ClusterRole
  name: terraform-runner
  apiGroup: rbac.authorization.k8s.io
```
6. If you're testing with private git repos, you need to create the known hosts config map
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: terraform-operator-known-hosts
data:
  known_hosts: |-
    bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
    github.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=
    github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl
    github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
    gitlab.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
    gitlab.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf
    gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9
    ssh.dev.azure.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
    vs-ssh.visualstudio.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
```
7. Install the manifests above and install the CRD. See [Dependencies](#dependencies)

# Building and Running the operator

## Basics
The scaffolding for the project is generated using `Kubebuilder`. It is a good idea to become familiar with this [project](https://github.com/kubernetes-sigs/kubebuilder). The [quick start](https://book.kubebuilder.io/quick-start.html) guide is also quite useful.

See `Makefile` at the root directory of the project. By default, executing `make` will build the project and produce an executable at `./bin/manager`

## Dependencies
To run successfully, any CRDs defined in the project should be regenerated and installed

The following steps should illustrate what is required before the project can be run:
1. `go mod tidy` - download the dependencies (this can take a while and there is no progress bar - need to be patient for this one)
2. `make manifests` - regenerates the CRD manifests
3. `make install` -  installs the CRDs into the cluster
4. `make generate` - generate the code

## Running the controller 
```
make run
```

## Running Tests
```
make test
```

## Extending the Library
As previously mentioned, familiarity with `kubebuilder` is required for developing this operator. Kubebuilder generates the scaffolding for new Kubernetes APIs.

```
$ kubebuilder create api --group run --version v1alpha1 --kind [YOUR_KIND]
 
Create Resource [y/n]
y
Create Controller [y/n]
y
```
Once you've developed your API, ensure to regenerate and install your CRDs. See [Dependencies](#dependencies)