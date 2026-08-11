# nginx-operator

`nginx-operator` is a Kubernetes Operator built with Kubebuilder.  
It manages a custom resource named `NginxCluster` (`platform.example.com/v1`) and reconciles each instance into:

- one `Deployment` running `nginx:1.27`
- one `Service` of type `ClusterIP`

The custom resource currently lets you set desired replica count via `spec.replicas`, and the operator reports observed readiness via `status.readyReplicas`.

## Quick Scan

For detailed architecture and design documentation, review the GitHub Wiki:
[Wiki Home](https://github.com/ninja-nb/nginx-operator/wiki).

If you want to quickly gauge platform/system design depth, start with the Wiki ADRs:

- [ADR-004 - Nutanix NKP Platform Architecture](https://github.com/ninja-nb/nginx-operator/wiki/01%20-%20ADR-004%20-%20Nutanix%20NKP%20Design)
- [ADR-001 - Fleet Application Deployment Platform (GitOps Pull Model)](https://github.com/ninja-nb/nginx-operator/wiki/02%20-%20ADR-001%20-%20Fleet%20Application%20Deployment%20Platform%20-%20System%20Design%20Wiki)
- [ADR-002 - Enterprise AI/ML Infrastructure on Kubernetes](https://github.com/ninja-nb/nginx-operator/wiki/03%20-%20ADR-002%20-%20Enterprise%20AI%20ML%20Infrastructure%20on%20Kubernetes)
- [ADR-003 - NGINX Operator Design Q&A](https://github.com/ninja-nb/nginx-operator/wiki/04%20-%20ADR-003%20-%20NGINX%20Operator%20Design%20Q%26A)

Then review implementation depth in this order:

1. `internal/controller/nginxcluster_controller.go` (reconciliation behavior)
2. `internal/controller/nginxcluster_controller_test.go` (test coverage and edge cases)
3. `api/v1/nginxcluster_types.go` (API design and validation markers)
4. `config/samples/platform_v1_nginxcluster.yaml` (runtime usage)

## How It Works

When an `NginxCluster` is created or updated, the controller:

1. reads the desired state (`spec.replicas`)
2. creates or updates a same-name `Deployment` in the same namespace
3. creates or updates a same-name `Service` exposing port `80`
4. sets owner references so child resources are garbage-collected with the CR
5. updates `status.readyReplicas` from the `Deployment` status

The controller watches:

- `NginxCluster` resources directly
- owned `Deployment` resources
- owned `Service` resources

## Repository Layout

- `api/v1/`: CRD Go types for `NginxCluster`
- `internal/controller/`: reconciliation logic
- `config/crd/`: generated CRD manifests
- `config/rbac/`: generated RBAC manifests
- `config/manager/`: controller manager deployment manifest
- `config/samples/`: sample custom resources
- `test/`: unit and e2e tests

## Documentation Placement (Recommended)

- Keep `README.md` as the **entry point**: project purpose, run/deploy commands, and quick links.
- Keep design narratives and interview/ADR content in the **GitHub Wiki** (best for browseable architecture storytelling).
- Keep operator implementation notes close to code under `docs/` only when they are code-adjacent and likely to evolve with commits.

## Prerequisites

- Go `1.24+`
- Docker
- `kubectl`
- Access to a Kubernetes cluster
- For e2e tests: [Kind](https://kind.sigs.k8s.io/)

## Quick Start (Local Development)

### 1) Run tests

```sh
make test
```

### 2) Run controller locally against your current kubeconfig context

```sh
make run
```

### 3) Apply sample custom resource

```sh
kubectl apply -f config/samples/platform_v1_nginxcluster.yaml
```

### 4) Verify resources

```sh
kubectl get nginxclusters.platform.example.com
kubectl get deploy,svc
```

## Build and Deploy to a Cluster

### Build and push an image

```sh
make docker-build docker-push IMG=<registry>/nginx-operator:<tag>
```

### Install CRDs

```sh
make install
```

### Deploy controller manager

```sh
make deploy IMG=<registry>/nginx-operator:<tag>
```

### Create a sample `NginxCluster`

```sh
kubectl apply -f config/samples/platform_v1_nginxcluster.yaml
```

## Uninstall

### Delete sample resources

```sh
kubectl delete -f config/samples/platform_v1_nginxcluster.yaml
```

### Undeploy controller manager

```sh
make undeploy
```

### Remove CRDs

```sh
make uninstall
```

## Useful Make Targets

- `make help`: list all targets
- `make manifests`: regenerate CRD/RBAC manifests
- `make generate`: regenerate DeepCopy code
- `make lint`: run linter
- `make lint-fix`: auto-fix lint issues
- `make test`: run unit/integration tests (non-e2e)
- `make test-e2e`: run e2e tests using a dedicated Kind cluster
- `make build-installer IMG=<registry>/nginx-operator:<tag>`: generate `dist/install.yaml`

## Distribution

### YAML bundle

Generate a single install manifest:

```sh
make build-installer IMG=<registry>/nginx-operator:<tag>
```

This creates `dist/install.yaml`, which can be published and installed with:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/<repo>/<ref>/dist/install.yaml
```

### Helm chart (optional)

Helm support can be scaffolded with:

```sh
kubebuilder edit --plugins=helm/v2-alpha
```

## Notes for Contributors

- Do not manually edit generated files under `config/crd/bases/`, `config/rbac/role.yaml`, or `zz_generated.*.go`.
- After changing API types (`api/v1/*_types.go`), run:

```sh
make manifests
make generate
```

- After changing Go logic, run:

```sh
make lint-fix
make test
```

## License

Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at:

http://www.apache.org/licenses/LICENSE-2.0

