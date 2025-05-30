package ru.itmo.wp.form;

import ru.itmo.wp.domain.Post;
import ru.itmo.wp.domain.User;

import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Pattern;
import javax.validation.constraints.Size;

public class CommentCredentials {
    @NotNull
    @NotEmpty
    @Size(min = 1, max = 3000)
    private String text;

    private User user;
    private Post post;

    public User getUser() {
        return user;
    }

    public Post getPost() { return post;}
    public void setUser(User user) {this.user = user;}

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }
}