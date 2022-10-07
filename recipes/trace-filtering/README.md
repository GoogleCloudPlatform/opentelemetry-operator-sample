# Filtering out spans

This recipe demonstrates how to configure the OpenTelemetry Collector
to filter out spans.

This recipe is based on applying a Collector config that includes the `filter` processor.
It provides an `OpenTelemetryCollector` object that, when created, instructs the Operator to
create a new instance of the Collector with that config. If overwriting an existing `OpenTelemetryCollector`
object (i.e., you already have a running Collector through the Operator such as the one from the
[main README](../../README.md#starting-the-collector)), the Operator will update that existing
Collector with the new config.


## Prerequisites

* A running Kubernetes cluster
* The OpenTelemetry Operator installed in your cluster
* A Collector deployed with the Operator (recommended)
* [One of the sample apps](../../sample-apps) from this repo installed in the cluster

Note that the `OpenTelemetryCollector` object needs to be in the same namespace as your sample
app, or the Collector endpoint needs to be updated to point to the correct service address.

## Running

### Deploying the Recipe

Apply the `OpenTelemetryCollector` object from this recipe:

```
kubectl apply -f collector-config.yaml
```

(This will overwrite any existing collector config, or create a new one if none exists.)

### Checking the filtered Spans

To stream logs from the otel-collector, run:
```
kubectl logs deployment/otel-collector -f
```

You should see only client spans, where `service.name` is `<sample>-app`. You should not see any spans from the server, where `service.name` is `<sample>-service`, as those are filtered out.

## Learn More

The filter processor can filter on more than just service name! For a complete description of its capabilities, see the upstream [README](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/filterprocessor#filter-processor).