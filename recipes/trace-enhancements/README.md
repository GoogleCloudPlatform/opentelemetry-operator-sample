# Trace enhancements

This recipe shows how to do basic trace enhancement by inserting values with the
[attributes processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/attributesprocessor).

## Prerequisites

* OpenTelemetry Operator installed in your cluster
* Running un-instrumented application (such as one of the [sample apps](../../sample-apps)).
* An `Instrumentation` object already created such as the one from the main [README](../../README.md#auto-instrumenting-applications)

# Running

Apply the `OpenTelemetryCollector` object from [`collector-config.yaml`](collector-config.yaml):

```
kubectl apply -f collector-config.yaml
```

# Combined with Resource Detection

Similar to the attributes processor, you can use the
[resource processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor)
to parse resource attributes.

Combine the resource processor  with the
[GCE/GKE resource detection recipe](../resource-detection)
to auto-populate a new attribute dynamically based on where the app is running.

The config file [`collector-config-resource-detection.yaml`](collector-config-resource-detection.yaml)
creates a new `location` attribute that is parsed from the cloud zone or region the pod is running in.

Create this config with:

```
kubectl apply -f collector-config-resource-detection.yaml
```

(Optionally, also enable [Cloud Trace integration](../cloud-trace) to see the traces in
your GCP dashboard.)

The regex in this config parses the first section of your node's zone or region to
its own `location` attribute, for example `us-central1` becomes `us` and `asia-south1-b`
becomes `asia`.
