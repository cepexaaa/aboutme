import React, {useRef, useState} from 'react';
import {useNavigate} from "react-router-dom";
import axios from "axios";

const WritePost = () => {

    const titleInputRef = useRef(null)
    const textInputRef = useRef(null)
    const router = useNavigate()
    const [error, setError] = useState('')

    const handleSubmit = async (event) => {
        event.preventDefault()
        const title = titleInputRef.current.value
        const text = textInputRef.current.value
        if (title.trim().length === 0 || text.trim().length === 0) {
            setError('Title or text can not be empty')
            return
        }
        try {
            const jwt = localStorage.getItem("jwt"); // Получаем JWT из localStorage
            const response = await axios.post("/api/posts", {
                title: title,
                text: text,
            }, {
                headers: {
                    Authorization: `Bearer ${jwt}` // Передаем JWT в заголовке
                }
            });
            if (response.data) {
                console.log(response)
                router("/");
            }
        } catch (error) {
            console.error('Error creating post:', error);
            setError('Failed to create post');
        }
        router("/")
    };

    return (
        <div className="form">
            <div className="header">Write Post</div>
            <div className="body">
                <form method="post" action="" onSubmit={handleSubmit}>
                    <input type="hidden" name="action" value="writePost"/>
                    <div className="field">
                        <div className="name">
                            <label htmlFor="title">Title</label>
                        </div>
                        <div className="value">
                            <input
                                autoFocus
                                id="title"
                                name="title"
                                ref={titleInputRef}
                                onChange={() => setError(null)}
                            />
                        </div>
                    </div>
                    <div className="field">
                        <div className="name">
                            <label htmlFor="text">Text</label>
                        </div>
                        <div className="value">
                            <textarea
                                id="text"
                                name="text"
                                ref={textInputRef}
                                onChange={() => setError(null)}
                            />
                        </div>
                    </div>
                    <div className="button-field">
                        <input type="submit" value="Write"/>
                    </div>
                    {error
                        ? <div className={'error'}>{error}</div>
                        : null
                    }
                </form>
            </div>
        </div>
    );
};

export default WritePost;