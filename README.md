# OpenTelemetry Operator Sample

This repo hosts samples for working with the OpenTelemetry Operator on GCP.

## Prerequisites

* A running GKE cluster
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

## Installing the OpenTelemetry Operator

Install the latest release of the Operator with:

```
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
```

## Starting the Collector

Set up an instance of the OpenTelemetry Collector by creating an `OpenTelemetryCollector` object.
The one in this repo sets up a basic OTLP receiver and logging exporter:

```
kubectl apply -f collector-config.yaml
```

## Auto-instrumenting Applications

The Operator offers [auto-instrumentation of application pods](https://github.com/open-telemetry/opentelemetry-operator#opentelemetry-auto-instrumentation-injection)
by adding an annotation to the Pod spec.

First, create an `Instrumentation` [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
that contains the settings for the instrumentation. We have provided a sample resource
in [`instrumentation.yaml`](instrumentation.yaml):

```
kubectl apply -f instrumentation.yaml
```


## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## License

Apache 2.0; see [`LICENSE`](LICENSE) for details.
