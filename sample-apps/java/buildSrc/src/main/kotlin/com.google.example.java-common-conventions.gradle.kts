plugins {
    // Apply the java Plugin to add support for Java.
    java
}

repositories {
    // Use Maven Central for resolving dependencies.
    mavenCentral()
}

dependencies {
    constraints {
        // Define dependency versions as constraints
        implementation("com.google.cloud:google-cloud-core:2.0.5")
        implementation("io.opentelemetry:opentelemetry-api:1.9.0")
        implementation("ch.qos.logback:logback-core:1.2.6")
        implementation("ch.qos.logback:logback-classic:1.2.2")
        // Allow j.u.l to pass through SLF4J into logback.
        implementation("org.slf4j:slf4j-api:1.7.32")
        implementation("org.slf4j:jul-to-slf4j:1.7.32")
    }

    // Use JUnit Jupiter for testing.
    testImplementation("org.junit.jupiter:junit-jupiter:5.7.2")
}

tasks.test {
    // Use JUnit Platform for unit tests.
    useJUnitPlatform()
}
