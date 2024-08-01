plugins {
    // Apply the common convention plugin for shared build configuration between library and application projects.
    id("com.google.example.java-common-conventions")

    // Apply the spring boot application plugin
    id("org.springframework.boot")
    id("io.spring.dependency-management")
}

dependencies {
    constraints {
        implementation("org.springframework.boot:spring-boot-starter-web:2.4.5")
    }
}
