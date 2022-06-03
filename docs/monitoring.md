---
layout: default
title: Monitoring
nav_order: 6
---

# Monitoring
The controller writes the following Prometheus metrics.

- `tfo_workflow_total`: The total number of submitted workflows/runs
- `tfo_workflow_status`: The current status of a Terraform workflow/run resource reconciliation
- `tfo_workflow_duration_seconds`: The duration in seconds of a Terraform workflow/run

*The metrics can be scraped from the controller's `/metrics` endpoint, the default metrics address port is set to `8080`*