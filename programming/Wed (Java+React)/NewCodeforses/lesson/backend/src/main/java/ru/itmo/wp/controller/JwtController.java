package ru.itmo.wp.controller;

import org.springframework.validation.BindingResult;
import org.springframework.web.bind.WebDataBinder;
import org.springframework.web.bind.annotation.*;
import ru.itmo.wp.domain.Post;
import ru.itmo.wp.domain.User;
import ru.itmo.wp.exception.ValidationException;
import ru.itmo.wp.form.UserCredentials;
import ru.itmo.wp.form.validator.UserCredentialsEnterValidator;
import ru.itmo.wp.service.JwtService;
import ru.itmo.wp.service.PostService;
import ru.itmo.wp.service.UserService;

import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;

@RestController
@RequestMapping("/api")
public class JwtController {

    private final UserService userService;
    private final JwtService jwtService;
//    private final PostService postService;

    public JwtController(UserService userService, JwtService jwtService,/* PostService postService,*/ UserCredentialsEnterValidator enterValidator) {
        this.userService = userService;
        this.jwtService = jwtService;
//        this.postService = postService;
        this.enterValidator = enterValidator;
    }

    private final UserCredentialsEnterValidator enterValidator;

    @InitBinder("userCredentials")
    public void initBinder(WebDataBinder webDataBinder){
        webDataBinder.addValidators(enterValidator);
    }

    @PostMapping("/jwt")
    public String create(@RequestBody @Valid UserCredentials userCredentials, BindingResult bindingResult){
        if (bindingResult.hasErrors()){
            throw new ValidationException(bindingResult);
        }
        User user = userService.findByLoginAndPassword(userCredentials.getLogin(), userCredentials.getPassword());
        return jwtService.create(user);
    }

    @GetMapping("/jwt")
    public User find(@RequestParam String jwt){
        return jwtService.find(jwt);
    }

//    @PostMapping("/posts")
//    public Post createPost(@RequestBody Post post, HttpServletRequest request) {
//        // Получаем JWT из заголовка запроса
//        String jwt = request.getHeader("Authorization");
//        if (jwt != null && jwt.startsWith("Bearer ")) {
//            jwt = jwt.substring(7);
//        }
//
//        User user = jwtService.find(jwt);
//        if (user == null) {
//            throw new RuntimeException("User not found");
//        }
//
//        post.setUser(user);
//
//        return postService.save(post);
//    }
}
