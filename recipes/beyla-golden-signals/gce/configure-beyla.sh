#!/bin/bash
#
# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# This script installs and configures beyla on to the current GCE VM instance, installing necessary dependencies.
# The script also configures the installed google-cloud-ops-agent so that it can collect metrics emitted by Beyla.

# Update dependencies
sudo apt update

# Install Java
echo "Installing Java 17..."
sudo apt install -y openjdk-17-jdk

# Download Beyla
echo "Downloading Beyla..."
BEYLA_V1_7_RELEASE=https://github.com/grafana/beyla/releases/download/v1.7.0/beyla-linux-amd64-v1.7.0.tar.gz
curl -Lo beyla.tar.gz $BEYLA_V1_7_RELEASE
mkdir -p beyla-installation/
tar -xzf beyla.tar.gz -C beyla-installation/
# Move beyla executable to /usr/local/bin
# /usr/local/bin is the path on a default GCE instance running Debian
sudo cp beyla-installation/beyla /usr/local/bin

# Configuring Beyla Environment variables when running in direct mode
echo  "Configuring Beyla..."
# Beyla configuration options
export BEYLA_OPEN_PORT=8080
export BEYLA_TRACE_PRINTER=text
export BEYLA_SERVICE_NAMESPACE="otel-beyla-gce-sample-service-ns"

# OpenTelemetry export settings
export OTEL_EXPORTER_OTLP_PROTOCOL="grpc"
export OTEL_EXPORTER_OTLP_ENDPOINT="http://0.0.0.0:4317"
export OTEL_SERVICE_NAME="otel-beyla-gce-sample-service"
echo "Beyla Configured..."

# Apply the custom configuration to installed Google Cloud Ops Agent
echo "Configuring the Google Cloud Ops Agent..."
sudo cp ./google-cloud-ops-agent/config.yaml /etc/google-cloud-ops-agent/config.yaml
sudo systemctl restart google-cloud-ops-agent
echo "Google Cloud Ops Agent restarted..."
