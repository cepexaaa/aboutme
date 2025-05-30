import './App.css';
import Enter from "./components/Middle/Main/Enter/Enter";
import Index from "./components/Middle/Main/Index/Index";
import React, {useEffect, useState} from "react";
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Application from "./Application";
import axios from "axios";
import WritePost from "./components/Middle/Main/WritePost/WritePost";
import Post from "./components/Middle/Main/Post/Post";
import OnePost from "./components/Middle/Main/Post/OnePost";
import User from "./components/Middle/Main/User/User";
import Register from "./components/Middle/Main/Register/Register";

function App() {

    const [login, setLogin] = useState(null)
    // const [users, setUsers] = useState(null)
    useEffect(() => {
        if (localStorage.getItem("jwt")){
            axios.get("/api/jwt", {
                params: {
                    jwt: localStorage.getItem("jwt")
                }
            }).then((response)=>{
                localStorage.setItem("login", response.data.login);
                setLogin(response.data.login)
            }).catch((error)=>{
                console.log(error)
            })

            // axios.get("/api/users", {
            //     params: {
            //         jwt: localStorage.getItem("jwt")
            //     }
            // }).then((response)=>{
            //     localStorage.setItem("users", response.data.user);
            //     setLogin(response.data.login)
            // }).catch((error)=>{
            //     console.log(error)
            // })
        }
    }, []);

    return (
        <div className="App">
            <BrowserRouter>
                <Routes>
                    <Route
                        index={true}
                        element={<Application setLogin={setLogin} login={login} page={<Index />}/>}
                    />
                    <Route
                        path={'/enter'}
                        element={<Application login={login} page={<Enter setLogin={setLogin}/>}/>}
                    />
                    <Route
                        path={'/register'}
                        element={<Application login={login} page={<Register setLogin={setLogin}/>}/>}
                    />
                    <Route
                        path={'/writePost'}
                        element={<Application login={login} setLogin={setLogin} page={<WritePost />}/>}
                    />
                    <Route
                        path={'/posts'}
                        element={<Application login={login} setLogin={setLogin} page={<Post />}/>}
                    />
                    <Route
                        exact path={'/posts/:id'}
                        element={<Application login={login} setLogin={setLogin} page={<OnePost />}/>}
                    />
                    <Route
                        path={'/users'}
                        element={<Application login={login} setLogin={setLogin} page={<User />}/>}
                    />
                    <Route
                        exact path={'/users/:id'}
                        element={<Application login={login} setLogin={setLogin} page={<User />}/>}
                    />
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default App;
