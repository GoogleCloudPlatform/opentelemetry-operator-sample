# Sample apps

This directory holds sample apps in various languages for working with
the operator and auto-instrumentation. See below to get started:

* [NodeJS](nodejs)
* [Java](java)
* [Python](python)
* DotNET
* Go

## Setup

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

## Troubleshooting

See [troubleshooting.md](troubleshooting.md) if you encounter issues with any of these apps.
