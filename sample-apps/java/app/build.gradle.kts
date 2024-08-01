/*
 * This file was generated by the Gradle 'init' task.
 */

plugins {
    id("com.google.example.java-application-conventions")
}

dependencies {
    implementation(project(":utilities"))
    implementation("ch.qos.logback:logback-classic")
    implementation("org.slf4j:slf4j-api")
    implementation("org.slf4j:jul-to-slf4j")
}

val mainClassName = "com.google.example.app.App"

application {
    // Define the main class for the application.
    mainClass.set(mainClassName)
}

/**
 * Create a Fat Jar with all dependencies for easier execution
 */
val fatJar = task("fatJar", type = Jar::class) {
    dependsOn("compileJava")
    archiveClassifier.set("standalone")
    duplicatesStrategy = DuplicatesStrategy.EXCLUDE
    manifest {
        attributes(Pair("Main-Class", mainClassName))
    }
    val sourcesMain = sourceSets.main.get()
    val contents = configurations.runtimeClasspath.get()
            .map { if (it.isDirectory) it else zipTree(it) } + sourcesMain.output
    from(contents)
    dependsOn(":utilities:jar", ":app:processResources")
}

tasks.build {
    dependsOn(fatJar)
}
