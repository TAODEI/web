package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Db 定义数据库
var Db *sql.DB

//User 定义结构体
type User struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

//Session 创建
type Session struct {
	SessionID string
	UserName  string
	UserID    int
}

func init() {
	//Db, err := sql.Open("mysql", "tao:12345678@tcp(localhost:2333)/users")
	var err error
	Db, err = sql.Open("mysql", "tao:12345678@/users")
	if err != nil {
		panic(err)
	}
}

//注册
func register(w http.ResponseWriter, r *http.Request) {
	var p User
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if p.Name == "" {
		http.Error(w, "用户名不能为空", 400)
		return
	}
	sqlstr := "select password from user where name = ?"
	_, err := Db.Exec(sqlstr, p.Name)
	if err != nil {
		http.Error(w, "该用户已注册", 400)
		panic(err)
	}
	sqlstr = "insert into user(name,nickname,password) value(?,?,?)"
	_, err = Db.Exec(sqlstr, p.Name, p.Nickname, p.Password)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	fmt.Fprintln(w)
	w.Write([]byte("注册成功"))
	return
}

//登录
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p User
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	sqlstr := "select password,id from user where name = ? "
	row := Db.QueryRow(sqlstr, p.Name)

	var password string
	var ID int
	row.Scan(&password, &ID)
	if password == "" {
		http.Error(w, "用户不存在", 404)
		return
	}
	//cookie与session关联
	cookie := http.Cookie{
		Name:     p.Name,
		Value:    p.Name + "666",
		HttpOnly: true,
	}
	if p.Password == password {
		s := Session{
			SessionID: p.Name + "666",
			UserName:  p.Name,
			UserID:    ID,
		}
		//发送cookie
		http.SetCookie(w, &cookie)
		sqlstr = "select username from sessions where sessionID =?"
		row := Db.QueryRow(sqlstr, p.Name+"666")
		var username string
		row.Scan(&username)
		if username == p.Name {
			http.Error(w, "该用户已登录", 400)
			return
		}
		//插入session
		sqlstr = "insert into sessions values(?,?,?)"
		_, err := Db.Exec(sqlstr, s.SessionID, s.UserName, s.UserID)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(200)
		fmt.Fprintln(w)
		w.Write([]byte("登录成功"))
		return
	}
	http.Error(w, "密码错误", 400)
	return

}

//退出登录
func quit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p User
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cookie, _ := r.Cookie(p.Name)
	if cookie == nil {
		http.Error(w, "请先登录", 400)
		return
	}
	if ok := verify(cookie.Value); ok == false {
		http.Error(w, "请先登录", 400)
		return
	}
	//登录成功
	sqlstr := "delete from sessions where username = ? "
	_, err := Db.Exec(sqlstr, p.Name)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	fmt.Fprintln(w)
	w.Write([]byte("退出登录成功"))
	return
}

//查询
func info(w http.ResponseWriter, r *http.Request) {
	//设置响应头
	w.Header().Set("Content-Type", "application/json")
	var p User
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cookie, _ := r.Cookie(p.Name)
	if cookie == nil {
		http.Error(w, "请先登录", 400)
		return
	}
	if ok := verify(cookie.Value); ok == false {
		http.Error(w, "请先登录", 400)
		return
	}
	sqlstr := "select name,nickname,password from user where name = ?"
	row := Db.QueryRow(sqlstr, p.Name)
	var v User
	err := row.Scan(&v.Name, &v.Nickname, &v.Password)
	if err == nil {
		json, _ := json.Marshal(v)
		w.WriteHeader(200)
		fmt.Fprintln(w)
		w.Write([]byte("查询成功"))
		w.Write(json)
		return
	}
}

//注销
func del(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p User
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cookie, _ := r.Cookie(p.Name)
	if cookie == nil {
		http.Error(w, "请先登录", 400)
		return
	}
	if ok := verify(cookie.Value); ok == false {
		http.Error(w, "请先登录", 400)
		return
	}
	//登录成功

	sqlstr := "delete from user where name = ? "
	_, err := Db.Exec(sqlstr, p.Name)
	if err != nil {
		panic(err)
	}
	sqlstr = "delete from sessions where sessionID = ? "
	_, err = Db.Exec(sqlstr, p.Name)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	fmt.Fprintln(w)
	w.Write([]byte("注销成功"))
	return
}

//修改
func change(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p User
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	cookie, _ := r.Cookie(p.Name)
	if cookie == nil {
		http.Error(w, "请先登录", 400)
		return
	}
	if ok := verify(cookie.Value); ok == false {
		http.Error(w, "请先登录", 400)
		return
	}
	//登录成功
	sqlstr := "update user set password=?,nickname=? where name=?"
	_, err := Db.Exec(sqlstr, p.Password, p.Nickname, p.Name)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	fmt.Fprintln(w)
	w.Write([]byte("修改成功"))
	return
}

//验证登录状态
func verify(value string) bool {
	sqlstr := "select sessionID,username,userID from sessions where sessionID = ?"
	_, err := Db.Exec(sqlstr, value)
	if err != nil {
		fmt.Println("..............................")
		fmt.Println(err)
		fmt.Println("..............................")
		return false
	}
	row := Db.QueryRow(sqlstr, value)
	var s Session
	row.Scan(&s.SessionID, &s.UserName, &s.UserID)
	return true
}

func main() {
	fmt.Println("running...")

	http.HandleFunc("/res", register)
	http.HandleFunc("/log", login)
	http.HandleFunc("/info", info)
	http.HandleFunc("/del", del)
	http.HandleFunc("/quit", quit)
	http.HandleFunc("/cha", change)

	http.ListenAndServe(":23333", nil)
	fmt.Println("over")
}

/*
creat table user(id int primary key auto_increment,name varchar(100) not null unique,password varchar(100) not null,nickname varchar(100))

create table sessions(sessionID varchar(100) PRIMARY KEY,
username varchar(100) not null,
userID int not null,foreign key(userID) references user(id))


func AddUser() error {
	sqlstr := "insert into users(name,password,nickname) value(?,?,?)"
	//预编译
	inStmt， err ：= Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("预编译错误")
		panic(err)
	}//

	//执行
	_, err ：= Db.Exec(sqlstr, "1", "2", "3")
	if err != nil {
		panic(err)
	}
	return nil
}*/
