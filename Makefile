# Copyright 2022 Google LLC
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

CONTAINER_REGISTRY=otel-operator
REGISTRY_LOCATION=us-central1
GCLOUD_PROJECT ?= $(shell gcloud config get project)

.PHONY: setup
setup:
	gcloud artifacts repositories create ${CONTAINER_REGISTRY} --repository-format=docker --location=${REGISTRY_LOCATION} --description="OpenTelemetry Operator sample registry"

.PHONY: sample-replace
sample-replace:
	sed -i "s/%GCLOUD_PROJECT%/${GCLOUD_PROJECT}/g" k8s/app.yaml
	sed -i "s/%CONTAINER_REGISTRY%/${CONTAINER_REGISTRY}/g" k8s/app.yaml
	sed -i "s/%REGISTRY_LOCATION%/${REGISTRY_LOCATION}/g" k8s/app.yaml
	sed -i "s/%GCLOUD_PROJECT%/${GCLOUD_PROJECT}/g" k8s/service.yaml
	sed -i "s/%CONTAINER_REGISTRY%/${CONTAINER_REGISTRY}/g" k8s/service.yaml
	sed -i "s/%REGISTRY_LOCATION%/${REGISTRY_LOCATION}/g" k8s/service.yaml
