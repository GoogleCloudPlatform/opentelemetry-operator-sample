# GCE/GKE Resource Detection

This recipe demonstrates GKE and GCE resource detection in a collector config,
deployed with the operator. These resource detectors add the following metadata:

GKE:
```
* cloud.provider ("gcp")
* cloud.platform ("gcp_gke")
* k8s.cluster.name (name of the GKE cluster)
```

GCE:
```
* cloud.platform ("gcp_compute_engine")
* cloud.account.id
* cloud.region
* cloud.availability_zone
* host.id
* host.image.id
* host.type
```

## Prerequisites

* OpenTelemetry Operator installed in your cluster
* Running un-instrumented application (such as one of the [sample apps](../../sample-apps)).
* An `Instrumentation` object already created such as the one from the main [README](../../README.md#auto-instrumenting-applications)

# Running

Apply the `OpenTelemetryCollector` object from [`collector-config.yaml`](collector-config.yaml):

```
kubectl apply -f collector-config.yaml
```

