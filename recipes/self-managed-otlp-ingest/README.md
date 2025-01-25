# Using Self Managed OTLP ingest with OpenTelemetry Operator

This recipe shows how to use self-managed OTLP ingest to export telemetry from an auto-instrumented application deployed on a GKE cluster.

The example demonstrates a scenario where a user has an application deployed on a GKE cluster and now needs to instrument it and export the generated telemetry to Google Cloud.

## Setting up the cluster with an application

> [!TIP]
> This section outlines the steps to create a cluster & deploy a sample application. Feel free to skip this section and move to [Instrumenting deployed applications](#instrumenting-deployed-applications) if you already have a sample app and cluster setup. \
> Alternatively, you can run a convenience script that executes the steps in this section to give you a cluster with a running application. You can execute this script from the root of this recipe: `./setup-application.sh`. 

### Creating a GKE cluster

Run the following commands to create a GKE autopilot cluster. Alternatively, you can create a cluster using the Google Cloud Platform UI.

```shell
# Setup required environment variables, modify based on your preference
export PROJECT_ID="<your-project-id>"
export CLUSTER_NAME="opentelemetry-autoinstrument-cluster"
export CLUSTER_REGION="us-central1"
```

The following `gcloud` command will begin the cluster creation process.
```shell
gcloud beta container --project "${PROJECT_ID}" clusters create-auto "${CLUSTER_NAME}" --region "${CLUSTER_REGION}"
```

Connect `kubectl` to the created cluster:
```shell
gcloud container clusters get-credentials ${CLUSTER_NAME} --region ${CLUSTER_REGION} --project ${PROJECT_ID}
```

### Deploy sample application (Java) on your cluster

Deploy an un-instrumented application on your GKE cluster. The application used here is the un-instrumented version of the [Java Instrumentation Quickstart](https://github.com/GoogleCloudPlatform/opentelemetry-operations-java/tree/main/examples/instrumentation-quickstart).

To build the sample app, from the root of this recipe:
```shell
# We pull the application from GitHub and build it 
git clone https://github.com/GoogleCloudPlatform/opentelemetry-operations-java.git
pushd opentelemetry-operations-java/examples/instrumentation-quickstart && \
DOCKER_BUILDKIT=1  docker build -f uninstrumented.Dockerfile -t java-quickstart . && \
popd && \
rm -rf opentelemetry-operations-java
```

Next, push the locally built image to Google Artifact Registry:
```shell
# Export environment variables for configuring Artifact Registry
export CONTAINER_REGISTRY=<your-registry-name>
export REGISTRY_LOCATION=us-central1

# If you do not have an Artifact repository, make one.
# Alternatively, you can use any existing artifact registry
gcloud artifacts repositories create ${CONTAINER_REGISTRY} --repository-format=docker --location=${REGISTRY_LOCATION} --description="Sample applications to auto-instrument using OTel operator"

# Tag and push the built application image to Google Artifact Registry
docker tag java-quickstart:latest ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/java-quickstart:latest
docker push ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/java-quickstart:latest

# Additionally, we build a docker image that can simulate traffic on the deployed application by sending requests to our application
pushd traffic && \
DOCKER_BUILDKIT=1  docker build -f hey.Dockerfile -t ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/hey:latest . && \
docker push ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/hey:latest && \
popd
```

Finally, deploy the pushed application & traffic generator in your cluster:
```shell
# This relies on the environment variables exported in previous steps
kubectl kustomize k8s | envsubst | kubectl apply -f -
```

The traffic generator simulates traffic on the deployed endpoint by hitting the `/multi` endpoint at 1 QPS.  For more information about the application, view the [application readme](uninstrumented-app/examples/instrumentation-quickstart/README.md).\

## Deploy a self-managed OpenTelemetry Collector on the cluster

This step deploys an [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) instance in the created cluster. 
Our application will then be configured to send telemetry to this collector instance and the collector will export the generated telemetry to Google Cloud.

> [!IMPORTANT]
> The following IAM permissions are required to be configured if you have enabled GKE Workload Identity Federation. This example uses GKE Autopilot, which has Workload Identity Federation enabled by default.
> It is necessary to configure Workload Identity to let the applications running in this cluster authenticate to Google Cloud APIs. The collector that will be deployed in this cluster will require these permissions

For clusters with GKE Workload Identity Federation enabled, grant the relevant IAM permissions:
```shell
# This is the number corresponding to your project ID, you can get this by visiting the Google Cloud Platform home page.
export PROJECT_NUMBER=<your-project-number>

gcloud projects add-iam-policy-binding projects/$PROJECT_ID \
    --role=roles/logging.logWriter \
    --member=principal://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/$PROJECT_ID.svc.id.goog/subject/ns/opentelemetry/sa/opentelemetry-collector \
    --condition=None
gcloud projects add-iam-policy-binding projects/$PROJECT_ID \
    --role=roles/monitoring.metricWriter \
    --member=principal://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/$PROJECT_ID.svc.id.goog/subject/ns/opentelemetry/sa/opentelemetry-collector \
    --condition=None
gcloud projects add-iam-policy-binding projects/$PROJECT_ID \
    --role=roles/cloudtrace.agent \
    --member=principal://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/$PROJECT_ID.svc.id.goog/subject/ns/opentelemetry/sa/opentelemetry-collector \
    --condition=None
```
*Note: The permissions are being granted to a service account in a given namespace within the cluster. The collector deployment will be configured to use the same service account and namespace.*

A self-managed OpenTelemetry Collector can be easily deployed by using scripts from [OTLP Kubernetes Ingest](https://github.com/GoogleCloudPlatform/otlp-k8s-ingest/tree/main) package:

```shell
kubectl kustomize https://github.com/GoogleCloudPlatform/otlp-k8s-ingest/k8s/base | envsubst | kubectl apply -f -
```

This script creates a properly configured Collector with relevant permissions & recommended settings. \
At this point, there is an OpenTelemetry Collector instance running in the cluster that is ready to receive & export telemetry.

## Instrumenting deployed applications

Now that we have an application deployed on a GKE cluster & an OpenTelemetry Collector instance configured to receive telemetry, we can instrument the application using the OpenTelemetry operator.

### Installing OpenTelemetry Operator

We use the OpenTelemetry Operator to inject auto-instrumentation into the apps running in our cluster.\
For GKE Autopilot, install `cert-manager` with the following [Helm](https://helm.sh) commands:
```shell
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

Install OpenTelemetry Operator with the following command:
```shell
# Install OpenTelemetry Operator
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/download/v0.116.0/opentelemetry-operator.yaml
```

### Injecting language specific auto-instrumentation

Using the OpenTelemetry Operator, we can inject language specific auto-instrumentation in applications deployed in the cluster.\
This is done by first creating the Instrumentation CRD (Custom Resource Definition) and then using annotations to apply the instrumentation to a specific Kubernetes resource.

```shell
# Create the Instrumentation CRD
kubectl apply -f instrumentation.yaml
```

Next, annotate the deployment in which the created instrumentation needs to be injected:

```shell
# Annotate the deployment to inject the created instrumentation
# This will trigger a rolling restart of the deployment
kubectl patch deployment.apps/quickstart-app -n default -p '{"spec":{"template":{"metadata":{"annotations":{"instrumentation.opentelemetry.io/inject-java": "default/sample-java-auto-instrumentation"}}}}}'

# This sample shows both instrumentation and application deployment in the default namespace.
# However, they could be deployed in separate namespaces as well.
```

## Viewing the result

After injecting the application, the application should start producing telemetry which can be viewed in Google Cloud.\
The instrumentation used in the sample is currently configured to generate only traces and the traces can be viewed on the GCP UI using the Trace Explorer.

## Cleanup

You can delete the cluster and the artifact registry created using this sample to avoid unnecessary charges:
```shell
# Delete the created cluster & the artifact registry
gcloud container clusters delete ${CLUSTER_NAME} --location=${CLUSTER_REGION}
gcloud artifacts repositories delete ${CONTAINER_REGISTRY} --location=${REGISTRY_LOCATION}
```
