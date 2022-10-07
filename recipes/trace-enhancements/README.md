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

