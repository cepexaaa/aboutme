import React, {useEffect, useState} from 'react';
import OnePost from "../Post/OnePost";
import {useNavigate} from "react-router-dom";
import axios from "axios";

const Index = () => {
    const [posts, setPosts] = useState(null)
    useEffect(() => {
        axios.get("/api/posts").then((response)=>{
            setPosts(response.data)
        }).catch((error)=>{
            console.log(error)
        })
    }, []);
    const sortedPosts = posts ? posts.sort((a, b) => b.id - a.id) : null;
    const router = useNavigate()   /// { posts, setSelectedPost}

    const handlePostClick = (post) => {
    //     setSelectedPost(post);
        router("/posts/" + post.id)
    };

    return (
        <div>
            {/*<p> lslslsl</p>*/}
            {sortedPosts ?
            <> {sortedPosts.map((post) => (
                <article key={post.id}>
                    <OnePost  inputId={post.id} />
                </article>
            ))} </>
                :
                <div>No some posts :-(</div>
            }
        </div>
    );
};

export default Index;