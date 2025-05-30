import './App.css';
import React, {useEffect, useState} from "react";
import Middle from "./components/Middle/Middle";
import Footer from "./components/Footer/Footer";
import Header from "./components/Header/Header";
import axios from "axios";

function Application({page, login, setLogin}) {

    const [posts, setPosts] = useState(null)

    useEffect(() => {
        axios.get("/api/posts").then((response)=>{
            setPosts(response.data)
        }).catch((error)=>{
            console.log(error)
        })
    }, []);


    return (
        <div>
            <Header setLogin={setLogin} login={login}/>
            <Middle
                posts={posts}
                page={page}
            />
            <Footer/>
        </div>
    );
}

export default Application;
/*
TODO: как сделать приём данных в пост запросах в backend при помощи Credentionals, а не передавать jwt? - COMPLETE!!!
TODO: вынести все общие данные выше и передавать их
TODO: посмотреть для чего нужен  JWT
TODO: куча валидации и валидаторов применить.
useMemo - это хук, который сохраняет результат вызова функции (первый аргумент) и пересчитывает его только при изменении зависимостей
useCallback - возвращает одну и туже ссылку на функцию, до тех пор, пока не изменится одна из зависимостей.
useEffect - для замены некоторых методов жизненного цикла классового компонента;

* */