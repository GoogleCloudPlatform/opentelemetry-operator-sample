# Trace Remote Sampling Configuration Recipe

This recipe shows how to use the Jaeger Remote Sampling
protocol to provide fine-grained control of trace sampling.


This recipe provides three main pieces of configuration:

- A ConfigMap that provides fine-grained trace sampling controls for the entire cluster.
- An OpenTelemetryCollector deployment that will serve the control protocol to instrumentation.
- An Instrumentation configuration that will leverage the OpenTelemetryCollector deployment to look up trace sampling information.

## Prerequisites

* OpenTelemetry Operator installed in your cluster
* Running un-instrumented application (such as one of the [sample apps](../../sample-apps)).

## Running

Apply the `ConfigMap` object from [`sampling-config.yaml`](sampling-config.yaml)

```
kubectl apply -f remote-sampling-config.yaml
```

This creates the configuration for trace sampling. At any point, we can modify the `remote-sampling-config.yaml` file and reperform this step to adjust sampling on the cluster in the future (see [Tuning](#tuning))


Apply the `OpenTelemetryCollector` object from [`collector-config.yaml`](collector-config.yaml)

```
kubectl apply -f collector-config.yaml
```

Next, create the `Instrumentation` object from [`instrumetnation.yaml`](instrumentation.yaml) that will use the remote sampling service:

```
kubectl apply -f instrumentation.yaml
```

This creates a new object named `instrumentation/trace-remote-sampling` in the current namespace.

Annotate your application pods to use these settings by editing the `instrumentation.opentelemetry.io`
annotation using one of the following commands:

> Note that if the app is a standalone Pod you can
>`kubectl annotate` directly on the Pod, but if it is owned by a Deployment or other replica controller
> you must patch the metadata of the Pod template.

* **NodeJS:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-nodejs="trace-remote-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-nodejs": "trace-remote-sampling"}}}}}'
  ```

* **Java:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-java="trace-remote-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java": "trace-remote-sampling"}}}}}'
  ```

* **Python:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-python="trace-remote-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-python": "trace-remote-sampling"}}}}}'
  ```

* **DotNET:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-python="trace-remote-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-dotnet": "trace-remote-sampling"}}}}}'
  ```

# Tuning

To tune sampling in the cluster, simply update the `remote-sampling-config.yaml` and reapply:

```
kubectl apply -f remote-sampling-config.yaml
```

The OpenTelemetryCollector deployment should pick up changes within a few seconds of the ConfigMap rollout, and clients will further pull in those changes within a few minutes.

This can help, e.g. when needing to increase the sampling rate of a service for better observability "on the fly" and turn it back down after collecting enough data.

The format of the remote sampling configuration is [documented here](https://www.jaegertracing.io/docs/1.28/sampling/#collector-sampling-configuration).

An example configuration which disables tracing prometheus metrics and health checks would look as follows:

```json
{
  "default_strategy": {
    "type": "probabilistic",
    "param": 0.5,
    "operation_strategies": [
      {
        "operation": "GET /health",
        "type": "probabilistic",
        "param": 0.0
      },
      {
        "operation": "GET /metrics",
        "type": "probabilistic",
        "param": 0.0
      }
    ]
  }
}
```