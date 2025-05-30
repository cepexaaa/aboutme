package ru.itmo.wp.controller;

import org.springframework.web.bind.annotation.*;
import ru.itmo.wp.domain.Comment;
import ru.itmo.wp.form.CommentCredentials;
import ru.itmo.wp.service.CommentService;

import java.util.List;

@RestController
@RequestMapping("/api")
public class CommentController {
    private final CommentService commentService;
    public CommentController(CommentService commentService) { this.commentService = commentService; }

    @GetMapping("/comments")
    public List<Comment> getComments() { return commentService.findAll(); }

    @PostMapping("/comments")
    public Comment addComment(@RequestBody CommentCredentials commentCredentials) { return commentService.createComment(commentCredentials);}
}
