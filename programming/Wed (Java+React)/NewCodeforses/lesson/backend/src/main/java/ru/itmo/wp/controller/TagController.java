package ru.itmo.wp.controller;

import org.springframework.web.bind.annotation.*;
import ru.itmo.wp.domain.Tag;
import ru.itmo.wp.service.TagService;

import java.util.List;

@RestController
@RequestMapping("/api")
public class TagController {
    private final TagService tagService;

    public TagController(TagService tagService) {this.tagService = tagService;}


    @GetMapping("/posts/{postId}/tags")
    public List<Tag> getTagsByPostId(@PathVariable long postId) {
        return tagService.findAllByPostId(postId);
    }

    @PostMapping("/posts/{postId}/tags")
    public List<Tag> saveTag(@PathVariable long postId, @RequestBody List<String> tags) {
        return tagService.addTagsToPost(postId, tags);
//        tagService.createTagByPost(postId, tags);
//        return tagService.save(tags);
    }
}
