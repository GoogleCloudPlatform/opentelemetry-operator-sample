plugins {
    // Apply the common convention plugin for shared build configuration between library and application projects.
    id("com.google.example.java-common-conventions")

    // Apply the application plugin to add support for building a CLI application in Java.
    application
    id("com.google.cloud.tools.jib")
}

jib {
    containerizingMode = "packaged"
    from.image = "gcr.io/distroless/java-debian10:11"
}