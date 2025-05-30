import React, {useRef, useState} from 'react';
import OnePost from "./OnePost";
import {useNavigate} from "react-router-dom";
const Post = ({ post, user, addComment, setSelectedPost}) => {
    const textInputRef = useRef(null)
    const [error, setError] = useState('')
    const router = useNavigate()

    const handleCommentSubmit = (event) => {
        event.preventDefault();
        const commentText = textInputRef.current.value
        if (commentText.trim() === '') return;

        addComment(post.id, commentText);
        textInputRef.current.value = '';
        router("/")
    };

    if (!post) {
        return (<div>Loading...</div>)
    }

    return (
        <div>
            <article>
                <OnePost  inputId={post.id}/>
            </article>

            {user && (
                <div className="form">
                    <div className="header">Add Comment</div>
                    <div className="body">
                        <form method="post" action="" onSubmit={handleCommentSubmit}>
                            <div className="field">
                                <div className="name">
                                    <label htmlFor="text">Comment</label>
                                </div>
                                <div className="value">
                                    <textarea id="text" name="text" ref={textInputRef} onChange={() => setError(null)}></textarea>
                                </div>
                            </div>
                            <div className="button-field">
                                <input type="submit" value="Add Comment" />
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {post.comments && post.comments.length > 0 ? (  
                post.comments.map((comment) => (
                    <p className="commentText">
                        {(comment.split("\n")).map((line) => (
                            <p>{line}</p>
                            ))}
                    </p>
                ))
            ) : (
                <div>No comments yet.</div>
            )}
        </div>
    );
};

export default Post;





// const newComment = {
//     id: post.comments.length + 1,
//     user: user,
//     text: commentText,
//     creationTime: new Date().toISOString()
// };



// <div className="comment" key={comment.id}>
//     <div className="information">By {comment.user?.login || 'Unknown'}, {comment.creationTime || 'Unknown'}</div>
//     {comment.text.split('\n').map((line, index) => (
//         <p key={index}>{line}</p>
//     ))}
// </div>