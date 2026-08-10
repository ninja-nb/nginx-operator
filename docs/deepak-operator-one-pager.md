# Deepak Operator Question - One Pager

Prepared for interview follow-up  
Author: Narendra Bedarkar

## What Deepak Was Testing

Deepak's question ("If a CR spec change takes 10 minutes, what will you look for?") tests:
- **Operator internals understanding** (not just Kubernetes basics)
- **Debugging depth** across control loop stages
- **Ability to separate root cause layers** (controller vs workload vs API server)

He was evaluating whether I can diagnose latency across the full operator lifecycle:
1. CR spec updated
2. Reconcile triggered
3. Controller mutates child resources
4. Workload rollout reaches readiness
5. CR status reflects observed state

---

## How This Maps to `nginx-operator`

In this project:
- `NginxCluster.spec.replicas` = desired state
- Reconcile updates:
  - `Deployment`
  - `Service`
- Reconcile writes: `NginxCluster.status.readyReplicas`

So a "10-minute delay" can be isolated by checking:
- Was reconcile triggered late?
- Were `CreateOrUpdate` calls slow/failing/retrying?
- Did Deployment update quickly but Pods become Ready slowly?
- Was status update delayed/conflicted?

---

## Direct Answer I Can Give in Interview

"I would debug this in layers. First, verify event-to-reconcile latency from controller logs. Next, time each reconcile step: CR fetch, Deployment reconcile, Service reconcile, and status update. Then compare desired vs actual state propagation: if `spec.replicas` updates Deployment quickly, the delay is likely rollout-side (scheduling pressure, image pull latency, readiness probe failures), not controller logic. Finally, verify status update reliability and add conditions/metrics/events so future delays are diagnosable quickly."

---

## Why This Answer Is Strong

- It is **operator-specific** and implementation-aware.
- It uses **measurable checkpoints** instead of generic troubleshooting.
- It demonstrates **production thinking**: observability, idempotency, and reliability.
- It directly ties theory to a real codebase (`nginx-operator`).

---

## 30-Second Version

"Deepak's question is really about locating latency in the operator control loop. In my `nginx-operator`, I would trace CR update to reconcile trigger, reconcile API mutations, pod readiness, and status propagation. That tells me exactly whether delay is controller-side or rollout-side. I would also add conditions, metrics, and events for faster diagnosis at scale."

