/*
 * This file was generated by the Gradle 'init' task.
 */

plugins {
    id("com.google.example.java-application-conventions")
}

dependencies {
    implementation(project(":utilities"))
    // TODO: maybe move spring things to conventions
    implementation("org.springframework.boot:spring-boot-starter-web:2.4.5")
}

application {
    // Define the main class for the application.
    mainClass.set("com.google.example.service.Main")
}

jib {
    container.ports = listOf("8080")
}
