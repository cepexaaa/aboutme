import React, {useEffect, useRef, useState} from 'react';
import {useNavigate} from "react-router-dom";
import axios from "axios";

const Register = () => {

    const loginInputRef = useRef(null)
    const nameInputRef = useRef(null)
    const passwordInputRef = useRef(null)
    const [error, setError] = useState(null)
    const router = useNavigate()

    const [users, setUsers] = useState(null)

    useEffect(() => {
        axios.get("/api/users").then((response)=>{
            setUsers(response.data)
        }).catch((error)=>{
            console.log(error)
        })
    }, []);

    const handleSubmit = async (event) => {
        event.preventDefault()
        const login = loginInputRef.current.value
        const password = passwordInputRef.current.value

        const loginRegex = /^[a-z]{3,16}$/;
        if (!loginRegex.test(login)) {
            setError('Login must be between 3 and 16 lowercase Latin letters');
            return;
        }

        const isLoginUnique = !users.some((user) => user.login === login);
        if (!isLoginUnique) {
            setError('This login is already taken');
            return;
        }

        if (password.isEmpty || 2 > password.length > 20) {
            setError('Wrong password');
        }

        axios.post("/api/users", {
            login: login,
            password: password,
        }).catch((error)=>{
            console.log(error)
        })

        router("/enter");
    };

    return (
        <div className="register form-box">
            <div className="header">Register</div>
            <div className="body">
                <form method="post" action="" onSubmit={handleSubmit}>
                    <input type="hidden" name="action" value="register"/>
                    <div className="field">
                        <div className="name">
                            <label htmlFor="login">Login</label>
                        </div>
                        <div className="value">
                            <input
                                autoFocus
                                name="login"
                                ref={loginInputRef}
                                onChange={() => setError(null)}
                            />
                        </div>
                    </div>
                    <div className="field">
                        <div className="name">
                            <label htmlFor="password">Password</label>
                        </div>
                        <div className="value">
                            <input
                                name="password"
                                type="password"
                                ref={passwordInputRef}
                                onChange={() => setError(null)}
                            />
                        </div>
                    </div>
                    {error
                        ? <div className={'error'}>{error}</div>
                        : null
                    }
                    <div className="button-field">
                        <input type="submit" value="Register"/>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Register;






// try {
//     const response = await axios.post("/api/users", {
//         login: login,
//         password: password,
//     });
//     if (response.data) {
//         console.log(response)
//         router("/");
//     }
// } catch (error) {
//     console.error('Error creating post:', error);
//     setError('Failed to create post');
// }


// if (name.length <= 1 || name.length >= 32) {
//     setError('Name must be between 1 and 32 characters');
//     return;
// }

/*

<div className="field">
                        <div className="name">
                            <label htmlFor="name">Name</label>
                        </div>
                        <div className="value">
                            <input
                                name="name"
                                type="name"
                                ref={nameInputRef}
                                onChange={() => setError(null)}
                            />
                        </div>
                    </div>

 createUser({
            login: login,
            // name: name,
            password: password
        });
 */