# Kubernetes cluster metrics

This recipe demonstrates how to configure the OpenTelemetry Collector
(as deployed by the Operator) with the
[`k8sclusterreceiver`](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/k8sclusterreceiver).

The receiver watches the Kubernetes API server and emits cluster-level metrics
similar to kube-state-metrics. This recipe exports those metrics through
Google Managed Service for Prometheus and logs entity events with the debug
exporter.

## Prerequisites

* Cloud Monitoring API enabled in your GCP project
* The `roles/monitoring.metricWriter`
  [IAM permissions](https://cloud.google.com/trace/docs/iam#roles) for your
  cluster's service account, or Workload Identity configured for the collector
  service account
* A running GKE cluster
* The OpenTelemetry Operator installed in your cluster

The `k8sclusterreceiver` watches cluster-wide resources. This recipe uses one
collector replica and a dedicated service account with read-only RBAC for the
Kubernetes resources the receiver observes.

## Running

Create the service account and RBAC resources:

```
kubectl apply -f rbac.yaml
```

Apply the `OpenTelemetryCollector` object from this recipe:

```
kubectl apply -f collector-config.yaml
```

After the collector starts, check its logs:

```
kubectl logs deployment/k8s-cluster-receiver-collector -f
```

## Workload Identity Setup

If you use Workload Identity, bind a Google service account with permission to
write metrics to the collector's Kubernetes service account.

```
export GCLOUD_PROJECT=<your GCP project ID>
gcloud iam service-accounts create otel-k8s-cluster-receiver --project=${GCLOUD_PROJECT}
```

```
gcloud projects add-iam-policy-binding $GCLOUD_PROJECT \
    --member "serviceAccount:otel-k8s-cluster-receiver@${GCLOUD_PROJECT}.iam.gserviceaccount.com" \
    --role "roles/monitoring.metricWriter"
```

```
gcloud iam service-accounts add-iam-policy-binding \
    "otel-k8s-cluster-receiver@${GCLOUD_PROJECT}.iam.gserviceaccount.com" \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${GCLOUD_PROJECT}.svc.id.goog[default/otel-k8s-cluster-receiver]"
```

```
kubectl annotate serviceaccount otel-k8s-cluster-receiver \
    iam.gke.io/gcp-service-account=otel-k8s-cluster-receiver@${GCLOUD_PROJECT}.iam.gserviceaccount.com
```

Restart the collector pod after adding the annotation.

## View your Metrics

Navigate to console.cloud.google.com/monitoring/metrics-explorer, and in the
"Select a metric" dropdown, search for "prometheus/k8s" to see Kubernetes
cluster metrics.

## Troubleshooting

### PermissionDenied from Google Managed Service for Prometheus

Check that the collector service account can use a Google service account with
the `roles/monitoring.metricWriter` role.

### Forbidden from the Kubernetes API

Check that `rbac.yaml` was applied in the same namespace as the
`OpenTelemetryCollector` object, and that the collector pod is using the
`otel-k8s-cluster-receiver` service account.
