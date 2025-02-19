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
* (For private clusters): set up the [firewall rules](#firewall-rules) for cert-manager

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

#### Firewall rules

By default, private GKE clusters may not allow the necessary ports for
cert-manager to work, resulting in an error like the following:

```
Error from server (InternalError): error when creating "collector-config.yaml": Internal error occurred: failed calling webhook "mopentelemetrycollector.kb.io": failed to call webhook: Post "https://opentelemetry-operator-webhook-service.opentelemetry-operator-system.svc:443/mutate-opentelemetry-io-v1alpha1-opentelemetrycollector?timeout=10s": context deadline exceeded
```

To fix this, [create a firewall
rule](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters#add_firewall_rules)
for your cluster with the following command:

```
gcloud compute firewall-rules create cert-manager-9443 \
  --source-ranges ${GKE_MASTER_CIDR} \
  --target-tags ${GKE_MASTER_TAG}  \
  --allow TCP:9443
```

`$GKE_MASTER_CIDR` and `$GKE_MASTER_TAG` can be found by following the steps in
the [firewall
docs](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters#add_firewall_rules)
listed above.

### Installing the OpenTelemetry Operator

Install the Operator with:

```
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/download/v0.116.0/opentelemetry-operator.yaml
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

With a Collector and auto-instrumentation set up, you can experiment with it using one of the [sample applications](sample-apps),
or skip right to the [recipes](recipes) if you already have an application running.

## Sample Applications

The [`sample-apps/`](sample-apps/) folder contains basic apps to demonstrate collecting traces with
the operator in various languages:

* [NodeJS](sample-apps/nodejs)
* [Java](sample-apps/java)
* [Python](sample-apps/python)
* DotNET (coming soon)
* [Go](sample-apps/go)
* [NodeJS + Java](sample-apps/nodejs-java)

Each of these sample apps works well with the [recipes](recipes) listed below.

## Recipes

The [`recipes`](recipes/) directory holds different sample use cases for working with the
operator and auto-instrumentation along with setup guides for each recipe. Currently, there are:

* [Trace sampling configuration](recipes/trace-sampling)
* [Trace remote sampling config](recipes/trace-remote-sampling)
* [Trace filtering](recipes/trace-filtering)
* [Trace enhancements](recipes/trace-enhancements)
* [Cloud Trace integration](recipes/cloud-trace)
* [Resource detection](recipes/resource-detection)
* [Daemonset and Deployment](recipes/daemonset-and-deployment)
* [eBPF HTTP Golden Signals with Beyla](recipes/beyla-golden-signals)
* [eBPF HTTP Service Graph with Beyla](recipes/beyla-service-graph)

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for details.

## License

Apache 2.0; see [`LICENSE`](LICENSE) for details.
