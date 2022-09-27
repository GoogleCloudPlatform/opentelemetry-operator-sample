# Trace Sampling Configuration Recipe

This recipe shows how to set basic trace sampling configuration on auto-instrumented applications.

## Prerequisites

* OpenTelemetry Operator installed in your cluster
* Running un-instrumented application (such as one of the [sample apps](../../sample-apps)).

## Running

Create the `Instrumentation` object from [`instrumentation.yaml`](instrumentation.yaml):

```
kubectl apply -f instrumentation.yaml
```

This creates a new object named `instrumentation/trace-sampling` in the current namespace.
Edit the `spec.sampler` section of this file to update settings for 
[trace sampling](https://opentelemetry.io/docs/reference/specification/trace/tracestate-probability-sampling/),
such as the [type of sampler](https://opentelemetry.io/docs/reference/specification/trace/sdk/#built-in-samplers)
and sampling ratio.

Annotate your application pods to use these settings by editing the `instrumentation.opentelemetry.io`
annotation using one of the following commands:

> Note that if the app is a standalone Pod you can
>`kubectl annotate` directly on the Pod, but if it is owned by a Deployment or other replica controller
> you must patch the metadata of the Pod template.

* **NodeJS:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-nodejs="trace-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-nodejs": "trace-sampling"}}}}}'
  ```

* **Java:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-java="trace-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java": "trace-sampling"}}}}}'
  ```

* **Python:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-python="trace-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-python": "trace-sampling"}}}}}'
  ```

* **DotNET:**

  Pod:
  ```
  kubectl annotate pod/<APP-NAME> instrumentation.opentelemetry.io/inject-python="trace-sampling"
  ```
  Deployment:
  ```
  kubectl patch deployment.apps/<APP-NAME> -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-dotnet": "trace-sampling"}}}}}'
  ```
