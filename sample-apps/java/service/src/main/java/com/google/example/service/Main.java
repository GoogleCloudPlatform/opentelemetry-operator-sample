/*
 * Copyright 2021 Google
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.google.example.service;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.net.http.HttpClient;
import java.net.http.HttpResponse;
import java.net.http.HttpResponse.BodyHandlers;
import java.util.logging.Logger;
import java.util.logging.Level;

@RestController
@SpringBootApplication
public class Main {
    private static Logger logger = Logger.getLogger("Main");

    public static void main(String[] args) throws IOException {
        com.google.example.utilities.Logging.initializeLogging();
        SpringApplication.run(Main.class, args);
    }

    @GetMapping("/")
    public String home() {
        String service = System.getenv("NEXT_SERVER");
        if (service != null) {
            return remoteHelloWorld(service);
        }
        return helloWorld();
    }

    private String remoteHelloWorld(String nextServer) {
        try {
            HttpClient httpClient = HttpClient.newHttpClient();
            HttpResponse<String> response =
            httpClient.send(java.net.http.HttpRequest.newBuilder()
                .uri(java.net.URI.create(nextServer))
                .timeout(java.time.Duration.ofMinutes(2))
                .GET()
                .build(), BodyHandlers.ofString());
            return response.body();
        } catch (Exception e) {
            logger.log(Level.SEVERE, "failed to connect to next server", e);
            return "{\"response\":\"Failure\"}";
        }
    }

    private String helloWorld() {
        return "{\"response\":\"Hello World\"}";
    }
}
