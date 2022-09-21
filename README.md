# OpenTelemetry Operator Sample

This repo hosts samples for working with the OpenTelemetry Operator on GCP.

* [Running the Operator](#running-the-operator)
   * [Prerequisites](#prerequisites)
   * [Installing the OpenTelemetry Operator](#installing-the-opentelemetry-operator)
   * [Starting the Collector](#starting-the-collector)
   * [Auto-instrumenting Applications](#auto-instrumenting-applications)
* [Sample Applications](#sample-applications)
   * [NodeJS](#nodejs)
   * [Java](#java)
   * [Python](#python)
   * [Go](#go)
* [Contributing](#contributing)
* [License](#license)

## Running the Operator

### Prerequisites

* A running GKE cluster
* [Helm](https://helm.sh) (for GKE Autopilot)
* `cert-manager` installed in your cluster
  * [Instructions here](https://cert-manager.io/docs/installation/)

For GKE Autopilot, install `cert-manager` with the following [Helm](https://helm.sh) commands:
```
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install \
  --create-namespace \
  --namespace cert-manager \
  --set installCRDs=true \
  --set global.leaderElection.namespace=cert-manager \
  --set extraArgs={--issuer-ambient-credentials=true} \
  cert-manager jetstack/cert-manager
```

### Installing the OpenTelemetry Operator

Install the latest release of the Operator with:

```
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
```

### Starting the Collector

Set up an instance of the OpenTelemetry Collector by creating an `OpenTelemetryCollector` object.
The one in this repo sets up a basic OTLP receiver and logging exporter:

```
kubectl apply -f collector-config.yaml
```

### Auto-instrumenting Applications

The Operator offers [auto-instrumentation of application pods](https://github.com/open-telemetry/opentelemetry-operator#opentelemetry-auto-instrumentation-injection)
by adding an annotation to the Pod spec.

First, create an `Instrumentation` [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
that contains the settings for the instrumentation. We have provided a sample resource
in [`instrumentation.yaml`](instrumentation.yaml):

```
kubectl apply -f instrumentation.yaml
```

## Sample Applications

The [`sample-apps/`](sample-apps/) folder contains basic apps to demonstrate collecting traces with
the operator in various languages.

Each sample app can be built as a container and deployed in a GKE cluster with a few
commands. First, run `make setup` to create an Artifact Registry if you don't already
have one:

```
export REGISTRY_LOCATION=us-central1
export CONTAINER_REGISTRY=otel-operator
make setup
```

The build and deploy commands for each app use the environment variables from above to run
the container image from your registry, along with `GCLOUD_PROJECT`.

Make sure to set the `GCLOUD_PROJECT` environment variable before running these:

```
export GCLOUD_PROJECT=my-gke-project
```

### NodeJS

1. Build the sample app:
   ```
   make sample-nodejs-build
   ```
   This command will also update the local [manifests](sample-apps/nodejs/k8s)
   to refer to your image location.

2. Push the local image to the Artifact Registry you created
   in the setup steps (if you did not create one, or are using an already created registry,
   make sure to set the `REGISTRY_LOCATION` and `CONTAINER_REGISTRY` variables):
   ```
   make sample-nodejs-push
   ```

3. Deploy the app in your cluster:
   ```
   kubectl apply -f sample-apps/nodejs/k8s/.
   ```
   If you want to run the sample app in a specific namespace, pass `-n <your-namespace>`.

4. Run the following commands to patch the `app` and `server` deployments for auto-instrumentation:
   ```
   kubectl patch deployment.apps/nodeshowcase-app -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-nodejs": "true"}}}}}'
   kubectl patch deployment.apps/nodeshowcase-service -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-nodejs": "true"}}}}}'
   ```

### Java

Coming soon

### Python

Coming soon

### Go

Coming soon

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## License

Apache 2.0; see [`LICENSE`](LICENSE) for details.
