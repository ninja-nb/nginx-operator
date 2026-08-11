# ADR-004: Nutanix NKP Platform Architecture

## Status
Proposed

## Context
Reference architecture for running Kubernetes clusters with a management-plane-first design and Cluster API lifecycle operations.

## Architecture Diagram

```mermaid
flowchart TB
  User["Platform Engineer"]

  subgraph MGMT["Management Cluster (NKP)"]
    UI["UI / API"]
    Auth["Auth (OIDC/LDAP/RBAC/SSO)"]
    Fleet["Fleet / Workspaces"]
    CAPI["CAPI + Controllers"]
    LCM["Lifecycle Mgmt"]
    GitOps["GitOps Controller (Flux/Argo)"]
    Policy["Policy & Security"]
    Obs["Observability Stack"]
  end

  subgraph INFRA["Infrastructure Providers"]
    CAPX["CAPX (Nutanix)"]
    CAPA["CAPA (AWS)"]
    CAPG["CAPG (GCP)"]
    CAPZ["CAPZ (Azure)"]
    Prism["Prism Central APIs"]
  end

  subgraph WK["Workload Clusters Fleet"]
    W1["Nutanix Cluster"]
    W2["AWS Cluster"]
    W3["GCP Cluster"]
    W4["Azure Cluster"]
  end

  User --> UI
  UI --> Auth
  UI --> Fleet
  Fleet --> LCM
  LCM --> CAPI
  CAPI --> CAPX
  CAPI --> CAPA
  CAPI --> CAPG
  CAPI --> CAPZ
  CAPX --> Prism

  GitOps --> W1
  GitOps --> W2
  GitOps --> W3
  GitOps --> W4

  Policy --> W1
  Policy --> W2
  Policy --> W3
  Policy --> W4

  Obs --> W1
  Obs --> W2
  Obs --> W3
  Obs --> W4

  W1 -. telemetry/status .-> Obs
  W2 -. telemetry/status .-> Obs
  W3 -. telemetry/status .-> Obs
  W4 -. telemetry/status .-> Obs

  classDef mgmt fill:#e8f0fe,stroke:#1a73e8,color:#0b3d91,stroke-width:1px;
  classDef infra fill:#fff4e5,stroke:#ef6c00,color:#7a3e00,stroke-width:1px;
  classDef workload fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20,stroke-width:1px;
  classDef actor fill:#f3e5f5,stroke:#8e24aa,color:#4a148c,stroke-width:1px;

  class User actor;
  class UI,Auth,Fleet,CAPI,LCM,GitOps,Policy,Obs mgmt;
  class CAPX,CAPA,CAPG,CAPZ,Prism infra;
  class W1,W2,W3,W4 workload;
```

## Propagation Model

### Policy and Security to Workload Clusters

- Management defines baseline policy bundles and guardrails
- GitOps applies policy manifests into each workload cluster
- Cluster-local policy engines enforce continuously
- Compliance and drift status is aggregated back to management

### Observability to Workload Clusters

- Management defines a standard telemetry package
- GitOps deploys telemetry agents/collectors to each workload cluster
- Cluster-local collectors ship metrics, logs, and traces to central backends
- Central dashboards and alerting provide fleet-wide visibility

## Design Notes

- Management-first control model for provisioning and day-2 operations
- CAPI as declarative lifecycle API for create/scale/upgrade/repair
- Provider-specific infrastructure adapters under a common control pattern
- Parallel workload subsystems with consistent policy and telemetry controls
