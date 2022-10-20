# NodeJS+Java sample app

The purpose of this sample app is to demonstrate basic use of the operator across multi-lingual
microservice applications. It combines the [NodeJS sample app](../nodejs) with the
[Java sample service](../java), where the NodeJS app periodically calls to the Java service.

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

1. Build and push the [Java sample app](../java):
   ```
   pushd ../java
   make build
   popd
   ```

2. Build and push the [NodeJS sample app](../nodejs):
   ```
   pushd ../nodejs
   make build
   make push
   popd
   ```

3. Update the [manifests](k8s) in this directory to point to your newly-pushed images:
   ```
   make sample-replace
   ```

4. Deploy the apps in your cluster:
   ```
   kubectl apply -f k8s/.
   ```

5. Run the following commands to patch the app CronJob and service Deployment for auto-instrumentation:
   ```
   kubectl patch deployment.apps/nodeshowcase-app -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-nodejs": "true"}}}}}'
   kubectl patch deployment.apps/javashowcase-service -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java": "true"}}}}}'
   ```
   These commands will use the `Instrumentation` created as part of the Prerequisites.

## View your Spans

To stream logs from the otel-collector, which will include spans from this sample application, run:
```
kubectl logs deployment/otel-collector -f
```

Alternatively, follow the [cloud-trace recipe](../../recipes/cloud-trace/) to view your spans in Google Cloud Trace.
