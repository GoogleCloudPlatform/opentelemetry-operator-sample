plugins {
    // Support convention plugins written in Kotlin. Convention plugins are build scripts in 'src/main' that automatically become available as plugins in the main build.
    `kotlin-dsl`
}
repositories {
    mavenCentral()
    gradlePluginPortal()
}

dependencies {
    implementation("gradle.plugin.com.google.cloud.tools:jib-gradle-plugin:3.1.4")
    // version 3.x of the spring boot plugin requires a minimum Java 17 version
    implementation("org.springframework.boot:spring-boot-gradle-plugin:2.7.14")
}
