#!/bin/bash
# Build the server and client sample apps
SAMPLE_APP_JAVA_DIR=../../../sample-apps/java
SERVICE_JAR_FILE=service.jar

# Build client and server Java apps
pushd $SAMPLE_APP_JAVA_DIR
./gradlew :service:build
./gradlew :app:build
popd

# Run the server application
java -jar $SAMPLE_APP_JAVA_DIR/service/build/libs/$SERVICE_JAR_FILE &
