import React, {useEffect, useState} from 'react';
import axios from "axios";

const User = () => {

    const [allUsers, setUsers] = useState(null)

    useEffect(() => {
        axios.get("/api/users").then((response)=>{
            setUsers(response.data)
        }).catch((error)=>{
            console.log(error)
        })
    }, []);

    return (
        <div className="users datatable">
            <div className="caption">Users</div>
            <table>
                <thead>
                <tr>
                    <th>Id</th>
                    <th>Login</th>
                    <th>Count of posts</th>
                </tr>
                </thead>
                <tbody>
                {!allUsers ? (
                    <tr className="noData">
                        <td colSpan="3">No data</td>
                    </tr>
                ) : (
                    allUsers.map((user) => (
                        <tr key={user.id}>
                            <td className="id">{user.id}</td>
                            <td className="login">{user.login}</td>
                            <td className="name">{user.posts}</td>
                        </tr>
                    ))
                )}
                </tbody>
            </table>
        </div>
    );
};

export default User;

/*
{(allUsers || allUsers.length === 0) ? (
 */

//for comments and posts
/*
<@c.post post=post/>
    <#if user??>
        <div class="form">
            <div class="header">Add Comment</div>
            <div class="body">
                <form method="post" action="/post/${post.id}">
                    <div class="field">
                        <div class="name">
                            <label for="text">Comment</label>
                        </div>
                        <div class="value">
<#--                            <textarea id="text" name="text" value="${commentForm.text!}"></textarea>-->
                            <textarea id="text" name="text"></textarea>
                        </div>
<#--                        <@c.error "commentForm.text"/>-->
                    </div>
                    <div class="button-field">
                        <input type="submit" value="Add Comment">
                    </div>
                </form>
            </div>
        </div>
    </#if>

    <#list post.comments as comment>
        <div class="comment">
            <div class="information">By ${comment.user.login}, ${comment.creationTime}</div>
            <#list comment.text?split("\n") as line>
                <p>${line}</p>
            </#list>
        </div>
    </#list>

    <#if post.tags?? && post.tags?size != 0>
        <div class="tags">
            Tags:
            <#list post.tags as tag>
                <span class="tag">${tag.name}</span><#if tag_has_next>, </#if>
            </#list>
        </div>
    </#if>




 for list of posts

 <article>
        <a class="title" href="/post/${post.id}">${post.title}</a>
        <div class="information">By ${post.user.login}, ${post.creationTime}</div>
        <div class="body">${post.text}</div>
        <ul class="attachment">
            <li>Announcement of <a href="#">Codeforces Round #510 (Div. 1)</a></li>
            <li>Announcement of <a href="#">Codeforces Round #510 (Div. 2)</a></li>
        </ul>
        <div class="footer">
            <div class="left">
                <img src="<@spring.url '/img/voteup.png'/>" title="Vote Up" alt="Vote Up"/>
                <span class="positive-score">+173</span>
                <img src="<@spring.url '/img/votedown.png'/>" title="Vote Down" alt="Vote Down"/>
            </div>
            <div class="right">
                <img src="<@spring.url '/img/date_16x16.png'/>" title="Publish Time" alt="Publish Time"/>
                ${post.creationTime}
                <img src="<@spring.url '/img/comments_16x16.png'/>" title="Comments" alt="Comments"/>
                <a href="#">68</a>
            </div>
        </div>
    </article>

*/