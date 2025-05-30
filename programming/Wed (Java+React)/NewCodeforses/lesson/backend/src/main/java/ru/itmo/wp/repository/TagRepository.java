package ru.itmo.wp.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import ru.itmo.wp.domain.Tag;

import javax.transaction.Transactional;
import java.util.List;

public interface TagRepository extends JpaRepository<Tag, Long> {
    Tag findById(long id);

    Tag findByName(String name);

    @Query(value = "SELECT t.* FROM tag t JOIN post_tag pt ON t.id = pt.tag_id WHERE pt.post_id = ?1", nativeQuery = true)
    List<Tag> findAllByPostId(long postId);

    @Modifying
    @Transactional
    @Query(value = "INSERT INTO post_tag (post_id, tag_id) VALUES (?1, ?2)", nativeQuery = true)
    void createTagByPost(long postId, long tagId);

//    @Modifying
//    @Transactional
//    @Query(value = "INSERT INTO tag (tag) VALUES (?1)", nativeQuery = true)
//    Tag save(Tag tag);
}