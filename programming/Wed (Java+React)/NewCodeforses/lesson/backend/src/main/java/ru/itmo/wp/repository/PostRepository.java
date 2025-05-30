package ru.itmo.wp.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import ru.itmo.wp.domain.Post;

import java.util.List;

public interface PostRepository extends JpaRepository<Post, Long> {
    List<Post> findAllByOrderByCreationTimeDesc();
    List<Post> findAllByUser_Id(Long id);
//    Post save(Post post);
    Post findById(long id);

//    @Query(value = "SELECT p FROM Post p LEFT JOIN FETCH p.tags WHERE p.id = :id")
//    Post findByIdWithTags(@Param("id") Long id);
}
