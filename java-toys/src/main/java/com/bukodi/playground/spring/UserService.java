package com.bukodi.playground.spring;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping( "/users")
public class UserService {

    @GetMapping("/{id}")
    public User getPerson(@PathVariable Long id) {
        return null;
    }

    @GetMapping("/hello")
    public String handle() {
        return "Hello WebFlux";
    }
}
