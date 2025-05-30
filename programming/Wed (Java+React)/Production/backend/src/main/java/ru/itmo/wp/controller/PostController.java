package ru.itmo.wp.controller;

import org.springframework.web.bind.annotation.*;
import ru.itmo.wp.domain.Post;
import ru.itmo.wp.domain.User;
import ru.itmo.wp.service.JwtService;
import ru.itmo.wp.service.PostService;
import ru.itmo.wp.service.UserService;

import javax.servlet.http.HttpServletRequest;
import java.util.List;

@RestController
@RequestMapping("/api")
public class PostController {

    private final PostService postService;
    private final JwtService jwtService;
    private final UserService userService;

    public PostController(PostService postService, JwtService jwtService, UserService userService) {
        this.postService = postService;
        this.jwtService = jwtService;
        this.userService = userService;
    }

    @GetMapping("/posts")
    public List<Post> getAllPosts() {
        return postService.findAll();
    }

    @PostMapping("/posts")
    public Post createPost(@RequestBody Post post, HttpServletRequest request) {
        // Получаем JWT из заголовка запроса
        String jwt = request.getHeader("Authorization");
        if (jwt != null && jwt.startsWith("Bearer ")) {
            jwt = jwt.substring(7); // Убираем "Bearer " из токена
        }
        User user = jwtService.find(jwt);
        if (user == null) {
            throw new RuntimeException("User not found");
        }
        post.setUser(user);
        return postService.save(post);
    }

    @GetMapping("/posts/{id}")
    public Post getPostById(@PathVariable Long id) {
        return postService.findById(id);
    }

    @GetMapping("myPosts/{id}")
    public List<Post> getMyPosts(@PathVariable Long id) {return postService.findAllByUser_Id(id);}
}
