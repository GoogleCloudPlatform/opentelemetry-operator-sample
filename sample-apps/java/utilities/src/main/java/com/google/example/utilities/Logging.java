package com.google.example.utilities;

/**
 * Utilities for configuring logging.
 */
public final class Logging {
    private Logging() {}

    /** Initialize the java.util.logging -> slf4j redirect so logback gets baked-in logs. */
    public static void initializeLogging() {
        org.slf4j.bridge.SLF4JBridgeHandler.removeHandlersForRootLogger();
        org.slf4j.bridge.SLF4JBridgeHandler.install();
    }
}
