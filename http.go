package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 定义全局map(User)
var users = make(map[string]User)

// 定义全局map(cookie)
var cookies = make(map[string]string)

// 定义用户结构体(注意大小写)
type User struct {
	Name     string `json:"name"` // binding:"required"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
type user struct {
	Name   string `json:"name"`
	Cookie string `json:"cookie"`
}
//POST
func register(w http.ResponseWriter, r *http.Request) {
	var p User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&p)
	if p.Name == "" {
		w.WriteHeader(400)
		w.Write([]byte("注册失败"))
		return
	}
	//name := r.FormValue("name")
	if _, ok := users[p.Name]; ok {
		w.WriteHeader(400)
		w.Write([]byte("该用户已注册或已注销"))
		return
	}
	users[p.Name] = p
	w.WriteHeader(200)
	w.Write([]byte("注册成功"))
	return
}

//GET
func info(w http.ResponseWriter, r *http.Request) {
	//设置响应头
	w.Header().Set("Content-Type", "application/json")
	var p user
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&p)
	name := p.Name
	//cookie := r.FormValue("cookie")
	//name := r.FormValue("name")
	if cookies[name] == p.Cookie {
		user := users[name]
		json, _ := json.Marshal(user)
		w.Write([]byte("查询成功"))
		w.Write(json)
		return
	} else {
		w.Write([]byte("请先登录"))
		return
	}
}

//GET
func login(w http.ResponseWriter, r *http.Request) {
	/*	var p User
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&p)*/

	w.Header().Set("Content-Type", "application/json")
	var p User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&p)
	//name := r.FormValue("name")
	//password := r.FormValue("password")
	name := p.Name
	if v := users[name]; v.Name == "" {
		w.Write([]byte("用户已注销"))
		return
	}
	if v, ok := users[name]; ok && name != "" {
		if p.Password == v.Password {
			w.Write([]byte("登陆成功"))
			cookie := "666" + name
			cookies[name] = cookie
		} else {
			w.Write([]byte("密码错误"))
		}
	} else {
		w.Write([]byte("用户不存在"))
	}
}

//GET
func del(w http.ResponseWriter, r *http.Request) {
	//获取user参数
	w.Header().Set("Content-Type", "application/json")
	var p user
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&p)
	name := p.Name
	//name := r.FormValue("name")
	//cookie := r.FormValue("cookie")
	if p.Cookie == cookies[name] {
		v := users[name]
		v.Name = ""
		v.Password = ""
		v.Nickname = ""
		users[name] = v
		cookies[name] = ""
		//delete(users, name)
		//delete(cookies, name)
		w.Write([]byte("注销成功"))
		return
	} else {
		w.Write([]byte("请先登陆"))
		return
	}
}

//GET
//这里不想再来个结构体  所以没用JSON
func change(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.FormValue("name")
	cookie := r.FormValue("cookie")
	newpsw := r.FormValue("newpsw")
	newnickname := r.FormValue("newnickname")
	if cookie == cookies[name] {
		v := users[name]
		v.Password = newpsw
		v.Nickname = newnickname
		users[name] = v
		w.Write([]byte("修改成功"))
		return
	} else {
		w.Write([]byte("请先登陆"))
		return
	}
}
func main() {
	fmt.Println("running...")
	http.HandleFunc("/userres", register)
	http.HandleFunc("/userinfo", info)
	http.HandleFunc("/userlog", login)
	http.HandleFunc("/userdel", del)
	http.HandleFunc("/usercha", change)

	http.ListenAndServe(":2333", nil)
}
