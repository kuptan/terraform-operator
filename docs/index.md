---
layout: default
title: Overview
nav_order: 1
---

# Overview
The Terraform Operator provides support to run Terraform modules in Kubernetes in a declarative way as a [Kubernetes manifest](https://kubernetes.io/docs/concepts/cluster-administration/manage-deployment/)

The motivation behind the Terraform Operator was to run Terraform modules in Kubernetes in a GitOps environment (with Flux or ArgoCD). You might say there are tools that are already built for Kubernetes like CrossPlane. The idea is that we already have the modules and they work just great for us, we just wanted the option to run in Kubernetes without switching to a completely new tool.

The Terraform Operator can help you with just that.

## Features Highlight

* Define your terraform flows in Kubernetes declaritvely
* Specify the module source and version with abillity to pull from private git repos
* Select the target workspace and define your backend configuration
* Define variables from Kubernetes secrets and configmaps
* Define variable files from Kubernetes secrets and configmaps
* Dependency management where one terraform run depends on another terraform run
* Module outputs can be written to a Kubernetes secret
* Define retry limit on your terraform run in case of failure