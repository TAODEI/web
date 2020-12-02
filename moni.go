package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strings"
	//"github.com/howeyc/gopass"
)

type a struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	requestUrl := "http://work.muxi-tech.xyz/api/v1.0/auth/login/"
	var user a
	/* 输入账号和密码
	fmt.Print("账号：")
	fmt.Scanf("%s",&user.Email)
	fmt.Print("密码：")

	password, err := gopass.GetPasswdMasked()

	if err != nil {
		panic(err)
	}

	data := url.Values{}
	data.Set("username", "864978550@qq.com")
	data.Set("password", "bHN0YW8xMTIxOTE3")
	payload := strings.NewReader(data.Encode())*/

	user = a{"864978550@qq.com", "bHN0YW8xMTIxOTE3"}

	buf, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		panic(err)
	}
	payload := strings.NewReader(string(buf))
	request, err := http.NewRequest("POST", requestUrl, payload)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Content-Length", "169")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Origin", "http://work.muxi-tech.xyz")
	request.Header.Add("Referer", "http://work.muxi-tech.xyz/")
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Status)
	fmt.Println(string(body))
}
