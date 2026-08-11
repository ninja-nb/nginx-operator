# nginx-operator Wiki Home

This folder mirrors the GitHub Wiki content in Git so updates can be reviewed and versioned through normal PR workflow.

## ADR Index (Reading Flow)

1. [ADR-004: Nutanix NKP Platform Architecture](01 - ADR-004 - Nutanix NKP Design)
2. [ADR-001: Fleet Application Deployment Platform (GitOps Pull Model)](02 - ADR-001 - Fleet Application Deployment Platform - System Design Wiki)
3. [ADR-002: Enterprise AI/ML Infrastructure on Kubernetes](03 - ADR-002 - Enterprise AI ML Infrastructure on Kubernetes)
4. [ADR-003: NGINX Operator Design Q&A](04 - ADR-003 - NGINX Operator Design Q&A)

## Recommended Reading Order (Interview Friendly)

1. **ADR-004** for NKP architecture and platform understanding
2. **ADR-001** for fleet-wide deployment strategy and scaling model
3. **ADR-002** for AI/ML GPU scheduling and resource sharing trade-offs
4. **ADR-003** for operator implementation/debugging depth

## What This Demonstrates

- Platform architecture thinking (management-plane-first design)
- Kubernetes operations at scale (fleet, GitOps, policy, observability)
- Practical operator engineering and troubleshooting
- Trade-off reasoning across performance, reliability, and cost

## Actual Operator Implementation

Use these links to review the real implementation in the repository:

- [Controller Reconcile Logic (`internal/controller/nginxcluster_controller.go`)](https://github.com/ninja-nb/nginx-operator/blob/main/internal/controller/nginxcluster_controller.go)
- [Controller Tests (`internal/controller/nginxcluster_controller_test.go`)](https://github.com/ninja-nb/nginx-operator/blob/main/internal/controller/nginxcluster_controller_test.go)
- [CRD API Types (`api/v1/nginxcluster_types.go`)](https://github.com/ninja-nb/nginx-operator/blob/main/api/v1/nginxcluster_types.go)
- [Sample Custom Resource (`config/samples/platform_v1_nginxcluster.yaml`)](https://github.com/ninja-nb/nginx-operator/blob/main/config/samples/platform_v1_nginxcluster.yaml)
- [Project README (run/deploy/verify)](https://github.com/ninja-nb/nginx-operator/blob/main/README.md)

## Sync to GitHub Wiki

See `wiki/SYNC.md` for copy/push workflow to `nginx-operator.wiki.git`.
