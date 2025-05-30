import React, {useEffect, useState} from 'react';
import voteup from "../../../../assets/img/voteup.png";
import votedown from "../../../../assets/img/votedown.png";
import date1616 from "../../../../assets/img/date_16x16.png";
import comments1616 from "../../../../assets/img/comments_16x16.png";
import {useNavigate, useParams} from "react-router-dom";
import axios from "axios";

const OnePost = ({inputId}) => {

    const router = useNavigate();
    const handlePostClick = (id) => {
        router(`/posts/${id}`)
    };

    const {id} = useParams();
    const resId = id ? id : inputId;
    const [post, setPost] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchPost = async () => {
            try {
                const response = await axios.get(`/api/posts/${resId}`);
                setPost(response.data);
                setLoading(false);
            } catch (error) {
                console.error('Error fetching post:', error);
                setError('Failed to load post');
                setLoading(false);
            }
        };

        fetchPost();
    }, [id, resId]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }


    return (
        <div>
            <a className="title" href="#" onClick={() => handlePostClick(post.id)}>{post.title}</a>
            <div className="information">By {post.user?.login || 'Unknown'}, {post.creationTime || 'Unknown'}</div>
            <div className="body">{post.text}</div>
            <ul className="attachment">
                <li>Announcement of <a href="#">Codeforces Round #510 (Div. 1)</a></li>
                <li>Announcement of <a href="#">Codeforces Round #510 (Div. 2)</a></li>
            </ul>
            <div className="footer">
                <div className="left">
                    <img src={voteup} title="Vote Up" alt="Vote Up"/>
                    <span className="positive-score">+173</span>
                    <img src={votedown} title="Vote Down" alt="Vote Down"/>
                </div>
                <div className="right">
                    <img src={date1616} title="Publish Time" alt="Publish Time"/>
                    {post.creationTime || 'Unknown'}
                    <img src={comments1616} title="Comments" alt="Comments"/>
                    <a href="#">{post.comments ? post.comments.length : " ?"}</a>
                </div>
            </div>
        </div>
    )
}

export default OnePost;

/* from index
<a className="title" href="#" onClick={() => handlePostClick(post)}>{post.title}</a>
<div className="information">By {post.user?.login || 'Unknown'}, {post.creationTime || 'Unknown'}</div>
<div className="body">{post.text}</div>
<ul className="attachment">
    <li>Announcement of <a href="#">Codeforces Round #510 (Div. 1)</a></li>
    <li>Announcement of <a href="#">Codeforces Round #510 (Div. 2)</a></li>
</ul>
<div className="footer">
    <div className="left">
        <img src={voteup} title="Vote Up" alt="Vote Up"/>
        <span className="positive-score">+173</span>
        <img src={votedown} title="Vote Down" alt="Vote Down"/>
    </div>
    <div className="right">
        <img src={date1616} title="Publish Time" alt="Publish Time"/>
        {post.creationTime || 'Unknown'}
        <img src={comments1616} title="Comments" alt="Comments"/>
        <a href="#">{post.commentsCount || 0}</a>
    </div>
</div>



<a className="title" href="#">{post.title}</a>
                <div className="information">By {post.user?.login || 'Unknown'}, {post.creationTime || 'Unknown'}</div>
                <div className="body">{post.text}</div>
                <ul className="attachment">
                    <li>Announcement of <a href="#">Codeforces Round #510 (Div. 1)</a></li>
                    <li>Announcement of <a href="#">Codeforces Round #510 (Div. 2)</a></li>
                </ul>
                <div className="footer">
                    <div className="left">
                        <img src={voteup} title="Vote Up" alt="Vote Up"/>
                        <span className="positive-score">+173</span>
                        <img src={votedown} title="Vote Down" alt="Vote Down"/>
                    </div>
                    <div className="right">
                        <img src={date1616} title="Publish Time" alt="Publish Time"/>
                        {post.creationTime || 'Unknown'}
                        <img src={comments1616} title="Comments" alt="Comments"/>
                        <a href="#">{post.commentsCount || 0}</a>
                    </div>
                </div>
 */