#!/bin/bash
# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

UNSET_WARNING="Environment variable not set, please set the environment variable"

# Verify necessary environment variables are set
echo "${PROJECT_ID:?${UNSET_WARNING}}"
echo "${CLUSTER_NAME:?${UNSET_WARNING}}"
echo "${CLUSTER_REGION:?${UNSET_WARNING}}"
echo "${CONTAINER_REGISTRY:?${UNSET_WARNING}}"
echo "${REGISTRY_LOCATION:?${UNSET_WARNING}}"

echo "ENVIRONMENT VARIABLES VERIFIED"

echo "CREATING CLUSTER WITH NAME ${CLUSTER_NAME} in ${CLUSTER_REGION}"
gcloud beta container --project "${PROJECT_ID}" clusters create-auto "${CLUSTER_NAME}" --region "${CLUSTER_REGION}"
echo "CLUSTER CREATED SUCCESSFULLY"

echo "PULLING SAMPLE APPLICATION REPOSITORY"
echo "BUILDING SAMPLE APPLICATION IMAGE"
git clone https://github.com/GoogleCloudPlatform/opentelemetry-operations-java.git
pushd opentelemetry-operations-java/examples/instrumentation-quickstart && \
DOCKER_BUILDKIT=1  docker build -f uninstrumented.Dockerfile -t java-quickstart . && \
popd && \
rm -rf opentelemetry-operations-java
echo "APPLICATION IMAGE BUILT"

echo "CREATING CLOUD ARTIFACT REGISTRY"
gcloud artifacts repositories create ${CONTAINER_REGISTRY} --repository-format=docker --location=${REGISTRY_LOCATION} --description="Sample applications to auto-instrument using OTel operator"
echo "CREATED ${CONTAINER_REGISTRY} in ${REGISTRY_LOCATION}"

echo "PUSHING THE SAMPLE APPLICATION IMAGE TO THE REGISTRY"
docker tag java-quickstart:latest ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/java-quickstart:latest
docker push ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/java-quickstart:latest
echo "APPLICATION IMAGE PUSHED TO ARTIFACT REGISTRY"

echo "BUILDING TRAFFIC SIMULATOR IMAGE"
pushd traffic && \
DOCKER_BUILDKIT=1  docker build -f hey.Dockerfile -t ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/hey:latest . && \
docker push ${REGISTRY_LOCATION}-docker.pkg.dev/${PROJECT_ID}/${CONTAINER_REGISTRY}/hey:latest && \
popd
echo "TRAFFIC SIMULATOR IMAGE BUILT & PUSHED TO ARTIFACT REGISTRY"

echo "DEPLOYING APPLICATION ON ${CLUSTER_NAME}"
kubectl kustomize k8s | envsubst | kubectl apply -f -
echo "SAMPLE APPLICATION DEPLOYED"
