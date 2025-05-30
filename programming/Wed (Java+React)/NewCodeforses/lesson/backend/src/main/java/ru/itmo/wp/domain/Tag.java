package ru.itmo.wp.domain;

import org.hibernate.validator.constraints.UniqueElements;

import javax.persistence.*;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.Set;

@Entity
@Table(indexes = @Index(columnList = "name", unique = true))
public class Tag {
    @Id
    @GeneratedValue
    private long id;

    @NotNull
    @NotEmpty
    private String name;

//    @ManyToMany(mappedBy = "tags", fetch = FetchType.LAZY)
//    private Set<Post> posts;

    public Tag() {}
    public Tag(@NotNull String name) {this.name = name;}

    public long getId() {
        return id;
    }

    public void setId(long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
//    public Set<Post> getPosts() {
//        return posts;
//    }
//
//    public void setPosts(Set<Post> posts) {
//        this.posts = posts;
//    }
}