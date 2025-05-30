package ru.itmo.wp.controller;

import org.springframework.web.bind.annotation.*;
import ru.itmo.wp.domain.User;
import ru.itmo.wp.form.UserCredentials;
import ru.itmo.wp.service.UserService;

import java.util.List;

@RestController
@RequestMapping("/api")
public class UserController {
    private final UserService userService;

    public UserController(UserService userService) {this.userService = userService;}

    @GetMapping("/users")
    public List<User> getAllUsers() {
        System.out.println("SSSSSSSSSSSize" + userService.findAll().size());
        return userService.findAll();}

    @PostMapping("/users")
    public User createUser(@RequestBody UserCredentials userCredentials) {return userService.register(userCredentials);}
//        public User createUser(@RequestBody User user) {return userService.register2(user);}
}
