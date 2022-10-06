import gradle.kotlin.dsl.accessors._5cbe0d75cf8ee73441f17a7e57de6851.jib
import org.gradle.api.tasks.Copy
import org.gradle.kotlin.dsl.*

plugins {
    // Apply the common convention plugin for shared build configuration between library and application projects.
    id("com.google.example.java-application-conventions")

    // Apply the application plugin to add support for building a CLI application in Java.
    application
    id("com.google.cloud.tools.jib")
}

val agent by configurations.creating
val agentOutputDir = layout.buildDirectory.dir("otelagent").forUseAtConfigurationTime().get()

tasks.register<Copy>("copyAgent") {
    from (agent) {
        rename("exporter-auto(.*).jar", "gcp_ext.jar")
        rename("opentelemetry-javaagent(.*).jar", "otel_agent.jar")
    }
    into(agentOutputDir)
}

// TODO: Figure out how to share this across all three tasks.
tasks.named("jib") {
    dependsOn("copyAgent")
}
tasks.named("jibDockerBuild") {
    dependsOn("copyAgent")
}
tasks.named("jibBuildTar") {
    dependsOn("copyAgent")
}


jib {
    container.jvmFlags = mutableListOf(
        // Use the downloaded java agent.
        "-javaagent:/otelagent/otel_agent.jar",
        // Export every 5 minutes
        "-Dotel.metric.export.interval=5m",
        // Use the GCP exporter extensions.
        "-Dotel.javaagent.extensions=/otelagent/gcp_ext.jar",
        // Configure auto instrumentation.
        "-Dotel.traces.exporter=google_cloud_trace",
        "-Dotel.metrics.exporter=google_cloud_monitoring")
    extraDirectories {
        paths {
            path {
                into = "/otelagent"
                setFrom(agentOutputDir.asFile.toPath())
            }
        }
    }
}

dependencies {
    agent("io.opentelemetry.javaagent:opentelemetry-javaagent:1.13.1")
    agent("com.google.cloud.opentelemetry:exporter-auto:0.21.0-alpha")
}
