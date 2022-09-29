# Cloud Trace integration

This recipe demonstrates how to configure the OpenTelemetry Collector
(as deployed by the Operator) to send trace data to GCP [Cloud Trace](https://cloud.google.com/trace).

This recipe is based on applying a Collector config that enables the Google Cloud exporter.
It provides an `OpenTelemetryCollector` object that, when created, instructs the Operator to
create a new instance of the Collector with that config. If overwriting an existing `OpenTelemetryCollector`
object (i.e., you already have a running Collector through the Operator such as the one from the
[main README](../../README.md#starting-the-collector)), the Operator will update that existing
Collector with the new config.


## Prerequisites

* Cloud Trace API enabled in your GCP project
* The `roles/cloudtrace.agent` [IAM permission](https://cloud.google.com/trace/docs/iam#roles)
  for your cluster's service account (or Workload Identity setup as shown below).
* A running GKE cluster
* The OpenTelemetry Operator installed in your cluster
* A Collector deployed with the Operator (recommended)
* An application already deployed that is either:
  * Instrumented to send traces to the Collector
  * Auto-instrumented by the Operator
  * [One of the sample apps](../../sample-apps) from this repo

Note that the `OpenTelemetryCollector` object needs to be in the same namespace as your sample
app, or the Collector endpoint needs to be updated to point to the correct service address.

## Running

### Workload Identity Setup

If you have Workload Identity enabled (on by default in GKE Autopilot), you'll need to set
up a service account with permission to write traces to Cloud Trace. You can do this with
the following commands:

```
export GCLOUD_PROJECT=<your GCP project ID>
gcloud iam service-accounts create otel-collector --project=${GCLOUD_PROJECT}
```

Then give that service account permission to write traces:

```
gcloud projects add-iam-policy-binding $GCLOUD_PROJECT \
    --member "serviceAccount:otel-collector@${GCLOUD_PROJECT}.iam.gserviceaccount.com" \
    --role "roles/cloudtrace.agent"
```

Then bind the GCP service account to the Kubernetes ServiceAccount that is used by the Collector
you deployed in the prerequisites (note: set `$COLLECTOR_NAMESPACE` to the namespace you installed
the Collector in):

```
export COLLECTOR_NAMESPACE=default
gcloud iam service-accounts add-iam-policy-binding "otel-collector@${GCLOUD_PROJECT}.iam.gserviceaccount.com" \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${GCLOUD_PROJECT}.svc.id.goog[${COLLECTOR_NAMESPACE}/otel-collector]"
```

Finally, annotate the Collector's ServiceAccount to use Workload Identity:

```
kubectl annotate serviceaccount otel-collector \
    --namespace $COLLECTOR_NAMESPACE \
    iam.gke.io/gcp-service-account=otel-collector@${GCLOUD_PROJECT}.iam.gserviceaccount.com
```

### Deploying the Recipe

Apply the `OpenTelemetryCollector` object from this recipe:

```
kubectl apply -f collector-config.yaml
```

(This will overwrite any existing collector config, or create a new one if none exists.)

Once the Collector restarts, you should see traces from your application
