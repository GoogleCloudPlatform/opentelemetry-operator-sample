# Python sample app

This is a sample app consisting of a basic server written in Python. The
server listens for requests which the client makes on a timed loop.

## Prerequisites

* OpenTelemetry Operator installed in your cluster
* Artifact Registry set up in your GCP project (see the 
[main README.md](../../README.md#sample-applications))
* An `Instrumentation` object already created in the current namespace
  (such as [the sample `instrumentation.yaml`](#auto-instrumenting-applications)
  from the main [README](../../README.md))

## Running

1. Build the sample app:
   ```
   make build
   ```
   This command will also update the local [manifests](k8s)
   to refer to your image location.

2. Push the local image to the Artifact Registry you created
   in the setup steps (if you did not create one, or are using an already created registry,
   make sure to set the `REGISTRY_LOCATION` and `CONTAINER_REGISTRY` variables):
   ```
   make push
   ```

3. Deploy the app in your cluster:
   ```
   kubectl apply -f k8s/.
   ```
   If you want to run the sample app in a specific namespace, pass `-n <your-namespace>`.

4. Run the following commands to patch the `app` and `server` deployments for auto-instrumentation:
   ```
   kubectl patch deployment.apps/pythonshowcase-app -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-nodejs": "true"}}}}}'
   kubectl patch deployment.apps/pythonshowcase-service -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-python": "true"}}}}}'
   ```
   These commands will use the `Instrumentation` created as part of the Prerequisites.
