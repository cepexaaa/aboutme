package ru.itmo.wp.service;

import org.springframework.stereotype.Service;
import ru.itmo.wp.domain.Post;
import ru.itmo.wp.domain.Tag;
import ru.itmo.wp.repository.PostRepository;
import ru.itmo.wp.repository.TagRepository;

import javax.transaction.Transactional;
import java.util.ArrayList;
import java.util.List;

@Service
public class TagService {
    private final TagRepository tagRepository;
    private final PostRepository postRepository;

    public TagService(TagRepository tagRepository, PostRepository postService) {
        this.tagRepository = tagRepository;
        this.postRepository = postService;
    }

    public Tag findByName(String name) {
        return tagRepository.findByName(name);
    }

    public List<Tag> findAll() { return tagRepository.findAll(); }

    public Tag findById(long id) {return tagRepository.findById(id);}

//    public Tag save(Tag tag) {
//        return tagRepository.save(tag);
//    }

    public List<Tag> findAllByPostId(long postId) {return tagRepository.findAllByPostId(postId);}

    public void createTagByPost(long postId, long tagId) {tagRepository.createTagByPost(postId, tagId);}

    @Transactional
    public List<Tag> addTagsToPost(long postId, List<String> tagNames) {
        List<Tag> allTags = new ArrayList<>();
        for (String tagName : tagNames) {
            Tag tag = findByName(tagName);
            if (tag == null) {
                Tag newTag = new Tag(tagName);
                tagRepository.save(newTag);
                tagRepository.createTagByPost(postId, newTag.getId());
                allTags.add(newTag);
            } else {
                tagRepository.createTagByPost(postId, tag.getId());
                allTags.add(tag);
            }
        }
        return allTags;
//        Post post = postRepository.findById(postId)
//                .orElseThrow(() -> new RuntimeException("Post not found"));
//
//        for (String tagName : tagNames) {
//            Tag tag = tagRepository.findByName(tagName)
//                    .orElseGet(() -> {
//                        Tag newTag = new Tag();
//                        newTag.setName(tagName);
//                        return tagRepository.save(newTag);
//                    });
//
//            post.getTags().add(tag);
//        }

//        postRepository.save(post);
    }
}