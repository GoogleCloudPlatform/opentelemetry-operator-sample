plugins {
    id("com.google.example.java-library-conventions")
}

dependencies {
    implementation("com.google.cloud:google-cloud-core")
    implementation("io.opentelemetry:opentelemetry-api")
    implementation("ch.qos.logback:logback-core")
    implementation("ch.qos.logback:logback-classic")
    implementation("org.slf4j:slf4j-api")
    implementation("org.slf4j:jul-to-slf4j")
}
