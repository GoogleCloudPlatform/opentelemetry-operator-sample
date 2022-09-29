# OpenTelemetry Operator Sample

This repo hosts samples for working with the OpenTelemetry Operator on GCP.

* [Running the Operator](#running-the-operator)
   * [Prerequisites](#prerequisites)
   * [Installing the OpenTelemetry Operator](#installing-the-opentelemetry-operator)
   * [Starting the Collector](#starting-the-collector)
   * [Auto-instrumenting Applications](#auto-instrumenting-applications)
* [Sample Applications](#sample-applications)
* [Recipes](#recipes)
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

Then run one of the language-based sample apps below:

* [NodeJS](sample-apps/nodejs)
* Java (coming soon)
* Python (coming soon)
* DotNET (coming soon)
* Go (coming soon)

## Recipes

The [`recipes`](recipes/) directory holds different sample use cases for working with the
operator and auto-instrumentation along with setup guides for each recipe. Currently there are:

* [Trace sampling configuration](recipes/trace-sampling)
* Trace filtering (coming soon)
* Trace enhancements (coming soon)
* [Cloud Trace integration](recipes/cloud-trace)

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## License

Apache 2.0; see [`LICENSE`](LICENSE) for details.
