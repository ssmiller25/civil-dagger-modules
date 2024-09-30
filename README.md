# Civo Dagger Modules

A collection of dagger modules to interact with [Civo](https://www.civo.com/).

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/civo/dagger-modules?quickstart=1)

## Civo Cluster

Interact with Civo Kubernetes clusters.

***List Clusters***

```bash
dagger call cluster-list --api-token env:CIVO_TOKEN --region lon1
```
