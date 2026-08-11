#!/usr/bin/env bash
set -euo pipefail

KIND_CLUSTER="${KIND_CLUSTER:-kind-nginx-demo}"
SAMPLE_MANIFEST="${SAMPLE_MANIFEST:-config/samples/platform_v1_nginxcluster.yaml}"

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Missing required command: $cmd" >&2
    exit 1
  fi
}

require_cmd kind
require_cmd kubectl
require_cmd make

echo "Deleting Kind cluster: ${KIND_CLUSTER} (if present)"
kind delete cluster --name "${KIND_CLUSTER}" || true

echo "Creating Kind cluster: ${KIND_CLUSTER}"
kind create cluster --name "${KIND_CLUSTER}"

echo "Switching kubectl context to kind-${KIND_CLUSTER}"
kubectl config use-context "kind-${KIND_CLUSTER}" >/dev/null

echo "Generating and installing CRDs"
make manifests
make install

echo "Applying sample NginxCluster resource: ${SAMPLE_MANIFEST}"
kubectl apply -f "${SAMPLE_MANIFEST}"

echo "Done. Current NginxCluster resources:"
kubectl get nginxclusters.platform.example.com -A
