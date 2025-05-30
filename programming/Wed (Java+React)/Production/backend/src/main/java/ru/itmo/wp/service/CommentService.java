package ru.itmo.wp.service;

import org.springframework.stereotype.Service;
import ru.itmo.wp.domain.Comment;
import ru.itmo.wp.form.CommentCredentials;
import ru.itmo.wp.repository.CommentRepository;

import java.util.List;

@Service
public class CommentService {
    private final CommentRepository commentRepository;

    CommentService(CommentRepository commentRepository) {this.commentRepository = commentRepository;}

    public Comment findById(long id) {return commentRepository.findById(id);}

    public List<Comment> findAll() { return commentRepository.findAll();}

    public boolean isTextVacant(String text) {
        return !text.isEmpty();
    }

    public Comment createComment(CommentCredentials commentCredentials) {
        Comment comment = new Comment();
        comment.setText(commentCredentials.getText());
        comment.setUser(commentCredentials.getUser());
        comment.setPost(commentCredentials.getPost());
        commentRepository.save(comment);
        return commentRepository.save(comment);
    }
}
