# Sidecar Resource Detection

This recipe builds on the basic [resource detection](../resource-detection)
example by running a collector sidecar next to each instrumented workload and
forwarding enriched traces to a central collector service.

The sample keeps the sidecars in the same namespace as their applications while
the central collector runs in a shared `otel-observability` namespace. This is a
useful shape when different teams own application namespaces, but telemetry
export is managed centrally.

## Prerequisites

* A running GKE cluster
* The OpenTelemetry Operator installed in your cluster
* Python and Java sample applications from this repo, built and pushed as
  described in [sample-apps/python](../../sample-apps/python) and
  [sample-apps/java](../../sample-apps/java)

## Running

Create the namespaces and central collector:

```
kubectl apply -f namespaces.yaml
kubectl apply -f central-collector.yaml
```

Create the sidecar collectors and instrumentation resources:

```
kubectl apply -f python-sidecar.yaml
kubectl apply -f java-sidecar.yaml
```

Deploy the Python app into `otel-python` and the Java app into `otel-java`:

```
kubectl apply -n otel-python -f ../../sample-apps/python/k8s/.
kubectl apply -n otel-java -f ../../sample-apps/java/k8s/.
```

Patch the workloads so the Operator injects both auto-instrumentation and the
collector sidecar:

```
kubectl patch -n otel-python deployment.apps/pythonshowcase-app \
  -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-python":"my-instrumentation","sidecar.opentelemetry.io/inject":"otel-sidecar"}}}}}'
kubectl patch -n otel-python deployment.apps/pythonshowcase-service \
  -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-python":"my-instrumentation","sidecar.opentelemetry.io/inject":"otel-sidecar"}}}}}'

kubectl patch -n otel-java deployment.apps/javashowcase-service \
  -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java":"my-instrumentation","sidecar.opentelemetry.io/inject":"otel-sidecar"}}}}}'
kubectl patch -n otel-java cronjob.batch/javashowcase-app \
  -p '{"spec":{"jobTemplate":{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java":"my-instrumentation","sidecar.opentelemetry.io/inject":"otel-sidecar"}}}}}}}'
```

## Checking the modified spans

Stream logs from the central collector:

```
kubectl logs -n otel-observability deployment/otel-central-collector -f
```

The exported spans should include GKE and GCE resource attributes detected by
the sidecar collectors, such as:

```
     -> cloud.provider: STRING(gcp)
     -> cloud.platform: STRING(gcp_kubernetes_engine)
     -> k8s.cluster.name: STRING(autopilot-cluster-1)
     -> k8s.namespace.name: STRING(otel-python)
```

The application pods should also include the injected collector container:

```
kubectl get pods -n otel-python -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.spec.containers[*].name}{"\n"}{end}'
kubectl get pods -n otel-java -o jsonpath='{range .items[*]}{.metadata.name}{" "}{.spec.containers[*].name}{"\n"}{end}'
```
