# Java sample app

This sample app runs a server "service" that listens for requests from a client "app"
running as a CronJob.

## Prerequisites

* OpenTelemetry Operator installed in your cluster
* Artifact Registry set up in your GCP project (see the 
[main README.md](../../README.md#sample-applications))
* An `OpenTelemetryCollector` object already created in the current namespace,
  such as [the sample `collector-config.yaml`](../../README.md#starting-the-Collector)
  from the main [README](../../README.md)
* An `Instrumentation` object already created in the current namespace,
  such as [the sample `instrumentation.yaml`](../../README.md#auto-instrumenting-applications)
  from the main [README](../../README.md)

## Running

1. Build the sample app and service:
   ```
   make build
   ```
   This command will also push the app's images to the Artifact Registry you set
   up in the Prerequisites. It also updates the local [manifests](k8s) to refer to
   your image registry.

2. Deploy the apps in your cluster:
   ```
   kubectl apply -f k8s/.
   ```
   If you want to run the sample app in a specific namespace, pass `-n <your-namespace>`.

3. Run the following commands to patch the app CronJob and service Deployment for auto-instrumentation:
   ```
   kubectl patch deployment.apps/javashowcase-service -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java": "true"}}}}}'
   kubectl patch cronjob.batch/javashowcase-app -p '{"spec":{"jobTemplate":{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java": "true"}}}}}}}'
   ```
   These commands will use the `Instrumentation` created as part of the Prerequisites.

## View your Spans

To stream logs from the otel-collector, which will include spans from this sample application, run:
```
kubectl logs deployment/otel-collector -f
```

Alternatively, follow the [cloud-trace recipe](../../recipes/cloud-trace/) to view your spans in Google Cloud Trace.
