# Kubernetes NGINX Operator — Session Handoff

## Goal

Build a Kubernetes Operator using **Kubebuilder** that manages NGINX at scale.

The desired demonstration is:

```
                    Kubernetes Cluster
┌─────────────────────────────────────────────────────┐
│                                                     │
│  Control Plane                                      │
│       │                                             │
│       ▼                                             │
│  Kubernetes API                                     │
│       ▲                                             │
│       │ watch/reconcile                             │
│       │                                             │
│  ┌──────────────────────┐                           │
│  │ NGINX Operator       │                           │
│  │                      │                           │
│  │ Reconcile()          │                           │
│  └──────────┬───────────┘                           │
│             │                                       │
│             ▼                                       │
│      NginxCluster CR                                │
│             │                                       │
│             ▼                                       │
│        Deployment                                   │
│          /    \                                     │
│         /      \                                    │
│     Worker 1   Worker 2                             │
│       NGINX       NGINX                             │
│       Pods        Pods                              │
│                                                     │
└─────────────────────────────────────────────────────┘

```

Eventually demonstrate changing:

```
spec:
  replicas: 3

```

to:

```
spec:
  replicas: 10

```

and show the Operator reconciling the Deployment.

---

# Environment

Mac Intel (`darwin/amd64`).

Docker:

```
Docker version 29.6.2

```

Go:

```
go version go1.26.5 darwin/amd64

```

Kubebuilder:

```
KubeBuilder v4.15.0
Kubernetes 1.36.0
Go 1.26.5

```

Homebrew had newer versions installed but old binaries were initially present in `/usr/local/bin`.

Old binaries:

```
/usr/local/bin/kubectl -> v1.17.9
/usr/local/bin/kind    -> v0.11.1

```

Homebrew versions:

```
/usr/local/opt/kubernetes-cli/bin/kubectl -> v1.36.3
/usr/local/opt/kind/bin/kind              -> v0.32.0

```

The old binaries were moved out of the way so the newer Homebrew versions could be used.

Current intended tooling:

```
kubectl      v1.36.3
kind         v0.32.0
go           v1.26.5
kubebuilder  v4.15.0

```

---

# Kind Cluster

Existing old cluster was a single-node cluster:

```
kind-control-plane
Kubernetes v1.21.1

```

We did NOT delete it.

A new test cluster was created specifically for the NGINX Operator:

```
Cluster name: nginx-demo
Context: kind-nginx-demo

```

It has:

```
nginx-demo-control-plane
nginx-demo-worker
nginx-demo-worker2

```

Current node output:

```
NAME                       STATUS   ROLES                  AGE     VERSION
nginx-demo-control-plane   Ready    control-plane,master   ...     v1.21.1
nginx-demo-worker          Ready    <none>                 ...     v1.21.1
nginx-demo-worker2         Ready    <none>                 ...     v1.21.1

```

The cluster is working.

Important: the nodes are Kubernetes **v1.21.1**, even though the client/Kubebuilder tools are much newer. We decided to proceed with the current cluster first rather than recreate it again.

Before continuing, verify:

```
kubectl config current-context

```

Expected:

```
kind-nginx-demo

```

and:

```
kubectl get nodes -o wide

```

---

# Operator Project

Project directory:

```
/Users/nbedarkar/Documents/NB-2026-Aug/kubernetes/nginx-operator

```

Current directory should be:

```
cd /Users/nbedarkar/Documents/NB-2026-Aug/kubernetes/nginx-operator

```

Project was initialized with:

```
kubebuilder init \
  --domain example.com \
  --repo github.com/narendra/nginx-operator

```

API was created with:

```
kubebuilder create api \
  --group platform \
  --version v1 \
  --kind NginxCluster

```

Answered:

```
Create Resource -> y
Create Controller -> y

```

Kubebuilder 4.15 generated the controller under:

```
internal/controller/nginxcluster_controller.go

```

NOT:

```
controllers/nginxcluster_controller.go

```

---

# Current API Definition

File:

```
api/v1/nginxcluster_types.go

```

The `NginxClusterSpec` was changed from the default `Foo` field to:

```
// NginxClusterSpec defines the desired state of NginxCluster.
type NginxClusterSpec struct {
	// Number of NGINX replicas.
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`
}

```

The status was changed to:

```
// NginxClusterStatus defines the observed state of NginxCluster.
type NginxClusterStatus struct {
	// Number of NGINX replicas that are ready.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

```

There was a temporary `controller-gen` marker error:

```
must apply listType to an array, found integer
must apply listMapKey to an array, found integer

```

This was fixed.

After the fix, all of these succeeded with exit code 0:

```
make generate
make manifests
make build

```

This is an important checkpoint: **the Kubebuilder project currently compiles successfully.**

---

# Generated Project Structure

Relevant files:

```
api/v1/groupversion_info.go
api/v1/nginxcluster_types.go
api/v1/zz_generated.deepcopy.go

internal/controller/nginxcluster_controller.go
internal/controller/nginxcluster_controller_test.go
internal/controller/suite_test.go

config/crd/
config/rbac/
config/samples/
config/manager/

cmd/main.go
Dockerfile
Makefile
go.mod
go.sum

```

---

# What Has NOT Been Done Yet

The actual NGINX reconciliation logic has **NOT yet been implemented**.

The current controller is still Kubebuilder-generated/default code.

We have NOT yet:

1. Implemented `Reconcile()`
2. Created an NGINX Deployment from the `NginxCluster`
3. Created an NGINX Service
4. Updated `NginxCluster.status.readyReplicas`
5. Tested reconciliation
6. Built the operator image
7. Loaded the operator image into kind
8. Deployed the operator into the cluster
9. Tested scaling from 3 to 10 replicas
10. Tested self-healing
11. Tested multi-worker distribution

---

# Immediate Next Steps

## Step 1 — Install the CRD

First verify:

```
kubectl config current-context

```

Should be:

```
kind-nginx-demo

```

Then:

```
make install

```

Verify:

```
kubectl get crd nginxclusters.platform.example.com

```

Then:

```
kubectl explain nginxcluster.spec

```

Expected to show:

```
replicas <integer>

```

---

## Step 2 — Inspect Current Controller

Before changing it, inspect:

```
cat internal/controller/nginxcluster_controller.go

```

The next session should look at this exact generated controller and then replace/adapt it.

Do NOT assume the old controller code from earlier messages is directly compatible with Kubebuilder 4.15.

---

# Intended Controller Behavior

The Operator should watch:

```
NginxCluster

```

and reconcile:

### 1. Deployment

Create/update:

```
Deployment/<nginxcluster-name>

```

with:

```
spec:
  replicas: <NginxCluster.spec.replicas>

```

NGINX container:

```
image: nginx:1.27
containerPort: 80

```

Add reasonable CPU/memory requests.

Add readiness probe:

```
GET /
port 80

```

### 2. Service

Create:

```
Service/<nginxcluster-name>

```

selecting:

```
app: <nginxcluster-name>

```

port:

```
80

```

### 3. Ownership

Set controller references so that the Operator owns:

```
Deployment
Service

```

This is important for demonstrating Kubernetes garbage collection and controller reconciliation.

### 4. Status

Update:

```
status:
  readyReplicas: 3

```

based on the Deployment status.

---

# Sample Custom Resource

Eventually create:

```
cat > nginxcluster.yaml <<'EOF'
apiVersion: platform.example.com/v1
kind: NginxCluster
metadata:
  name: demo-nginx
spec:
  replicas: 3
EOF

```

Then:

```
kubectl apply -f nginxcluster.yaml

```

Expected resources:

```
kubectl get nginxcluster
kubectl get deployment
kubectl get pods -o wide
kubectl get service

```

Expected architecture:

```
NginxCluster/demo-nginx
          |
          v
     NGINX Operator
          |
          +------> Deployment/demo-nginx
          |             |
          |             +-- NGINX pod
          |             +-- NGINX pod
          |             +-- NGINX pod
          |
          +------> Service/demo-nginx

```

---

# Scale Test

After the basic implementation works:

```
kubectl edit nginxcluster demo-nginx

```

or:

```
kubectl patch nginxcluster demo-nginx \
  --type merge \
  -p '{"spec":{"replicas":10}}'

```

Then:

```
kubectl get pods -o wide

```

The Deployment should move from:

```
3 replicas

```

to:

```
10 replicas

```

without directly modifying the Deployment.

This is the key demonstration of the Operator pattern:

```
User changes desired state
          |
          v
NginxCluster.spec.replicas = 10
          |
          v
Operator observes change
          |
          v
Reconcile()
          |
          v
Deployment replicas = 10
          |
          v
Kubernetes creates pods

```

---

# Failure / Self-Healing Test

After that, delete an NGINX pod:

```
kubectl delete pod <nginx-pod-name>

```

Then:

```
kubectl get pods -w

```

Kubernetes Deployment should recreate it.

This also demonstrates an important architectural distinction:

```
Operator
    |
    | manages desired Deployment state
    v
Deployment Controller
    |
    | manages desired Pod count
    v
Pods

```

The Operator doesn't need to create every Pod itself.

---

# Later: Run Operator Inside Kubernetes

Initially we can run:

```
make run

```

locally to simplify debugging.

After reconciliation works, build an image:

```
make docker-build IMG=nginx-operator:dev

```

Load it into kind:

```
kind load docker-image nginx-operator:dev --name nginx-demo

```

Then deploy:

```
make deploy IMG=nginx-operator:dev

```

Verify:

```
kubectl get pods -n nginx-operator-system

```

The exact namespace/name should be confirmed from the generated Kubebuilder manifests rather than assumed.

---

# Final Demonstration We Want

The eventual interview/demo story should be:

```
                    User
                     |
                     | kubectl apply
                     v
             NginxCluster CR
                     |
                     v
             Kubernetes API
                     |
                     v
             NGINX Operator
                     |
              Reconcile()
                     |
          +----------+----------+
          |                     |
          v                     v
     Deployment              Service
          |
          v
   +------+------+
   |             |
   v             v
Worker 1      Worker 2
   |             |
 NGINX         NGINX
 Pods          Pods

```

Then demonstrate:

```
replicas: 3
     |
     | change CR
     v
replicas: 10
     |
     v
Operator reconciles
     |
     v
Deployment = 10
     |
     v
10 NGINX Pods

```

This gives a concrete demonstration of:

- Kubernetes CRD
- Operator pattern
- Controller reconciliation
- Desired vs observed state
- Deployment ownership
- Service creation
- Status updates
- Kubernetes self-healing
- Horizontal scaling
- Multi-node scheduling
- Operator lifecycle
- Eventually operator HA/scaling