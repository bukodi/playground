package com.bukodi.playground.spring;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.reactive.server.WebTestClient;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class MyTests {

    @Autowired
    WebTestClient client;

//    @BeforeEach
//    void setUp(ApplicationContext context) {
//        client = WebTestClient.bindToApplicationContext(context).build();
//    }

    @Test
    void testHello() throws Exception {
        client.get().uri("/users/hello")
                .accept(MediaType.APPLICATION_JSON)
                .exchange()
                .expectStatus().isOk()
                .expectHeader().contentTypeCompatibleWith(MediaType.APPLICATION_JSON);
    }
}