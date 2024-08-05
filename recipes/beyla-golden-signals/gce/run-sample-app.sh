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
# Build the server and client sample apps
SAMPLE_APP_JAVA_DIR=../../../sample-apps/java
SERVICE_JAR_FILE=service.jar

# Build client and server Java apps
pushd $SAMPLE_APP_JAVA_DIR
./gradlew :service:build
./gradlew :app:build
popd

# Run the server application
java -jar $SAMPLE_APP_JAVA_DIR/service/build/libs/$SERVICE_JAR_FILE 2>&1 > logs_service.txt &
