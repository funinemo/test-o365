// https://outlook.office.com/のAPIを利用する場合、scopeでもhttps://outlook.office.com/をつける必要がある。
// https://msdn.microsoft.com/ja-jp/office/office365/api/mail-rest-operations
//   最低限必要なスコープ: 次のいずれか:
//     https://outlook.office.com/mail.read
//     wl.imap
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	//	"time"
)

type ConfigData struct {
	Client_id     string
	Client_secret string
	Redirect_uri  string
}

type TokenBody struct {
	Token_type     string
	Scope          string
	Access_token   string
	Refresh_token  string
	Expires_in     int
	ext_expires_in int
}

func auth_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<html lang="en">
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>"
		</head>
		<body><pre>`)
	log.Println("=====/auth START=====")
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	code := r.FormValue("code")
	state := r.FormValue("state")
	log.Printf("%#v", r.Form)
	log.Printf("code:%s\n", code)
	log.Printf("state:%s\n", state)

	for k, v := range r.Form {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
	var config ConfigData
	read_config(&config)
	log.Println("=====/auth -> token get START=====")
	urls := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("scope", "user.read mail.read")
	values.Add("client_id", config.Client_id)
	values.Add("redirect_uri", config.Redirect_uri)
	values.Add("client_secret", config.Client_secret)

	req, _ := http.NewRequest("POST", urls, strings.NewReader(values.Encode()))
	log.Printf("req1 : %s\n", req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("req set header : %s\n", dump)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dumpResp, _ := httputil.DumpResponse(resp, true)
	dumpRespBody, _ := ioutil.ReadAll(resp.Body)
	log.Printf("resp : %s\n", dumpResp)
	var tokens TokenBody

	if err := json.Unmarshal(dumpRespBody, &tokens); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Access_token: %s\n", tokens.Access_token)
	fmt.Fprintf(w, "Scope: %s\n", tokens.Scope)
	fmt.Fprintf(w, "Token_type: %s\n", tokens.Token_type)
	fmt.Fprintf(w, "Expires_in: %s\n", tokens.Expires_in)

	log.Printf("=====/auth -> Graph START=====")
	url := "https://graph.microsoft.com/v1.0/me"
	//url := "https://outlook.office.com/api/v2.0/me"
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", tokens.Token_type+" "+tokens.Access_token)
	dump, _ = httputil.DumpRequestOut(req, true)
	log.Printf("request header : %s\n", dump)
	resp, err = client.Do(req)
	dumpResp, _ = httputil.DumpResponse(resp, true)
	//	dumpRespBody, _ = ioutil.ReadAll(resp.Body)
	log.Printf("resp : %s\n", dumpResp)
	//
	//	time.Sleep(5 * time.Second)
	log.Printf("=====/auth -> MailGet START=====")
	//		url = "https://outlook.office.com/api/v2.0/me/events"
	//	url = "https://outlook.office.com/api/v2.0/me/mailfolders/inbox/messages?$top=5"
	url = "https://graph.microsoft.com/v1.0/me/mailfolders/inbox/messages?$top=1"
	//	url = "https://graph.microsoft.com/v1.0/me/events"
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", tokens.Token_type+" "+tokens.Access_token)
	dump, _ = httputil.DumpRequestOut(req, true)
	log.Printf("request header : %s\n", dump)
	resp, err = client.Do(req)
	dumpResp, _ = httputil.DumpResponse(resp, true)
	dumpRespBody, _ = ioutil.ReadAll(resp.Body)
	log.Printf("resp : %s\n", dumpResp)
	var buf map[string]interface{}
	err = json.Unmarshal(dumpRespBody, &buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================")
	log.Printf("%v\n", buf)

	var dat map[string]interface{}
	if err = json.Unmarshal(dumpRespBody, &dat); err != nil {
		log.Fatal(err)
	}
	log.Println("=========================")
	for k, v := range dat {
		fmt.Fprintf(w, "%v,%s", k, v)
	}
	fmt.Fprintln(w, `</pre></body></html>`)

}

// 試行錯誤中
func start_handler(w http.ResponseWriter, r *http.Request) {
	log.Println("=====/start START=====")
	var config ConfigData
	read_config(&config)
	log.Println("=====/auth -> token get START=====")
	urls := "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	values := url.Values{}
	values.Set("response_type", "code")
	values.Add("scope", "user.read mail.read")
	values.Add("client_id", config.Client_id)
	values.Add("state", "1234567890")
	values.Add("response_mode", "query")
	values.Add("redirect_uri", "http://localhost:8080/auth")

	req, _ := http.NewRequest("GET", urls, strings.NewReader(values.Encode()))
	log.Printf("req1 : %s\n", req)
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("req set header : %s\n", dump)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dumpResp, _ := httputil.DumpResponse(resp, true)
	//	dumpRespBody, _ := ioutil.ReadAll(resp.Body)
	log.Printf("resp : %s\n", dumpResp)
	fmt.Fprintf(w, "%s", dumpResp)

}
func root_handler(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Lmicroseconds)
	log.Println("=====/ START=====")
	log.SetFlags(log.Lshortfile)
	t := template.Must(template.ParseFiles("test.html"))
	err := t.ExecuteTemplate(w, "test.html", nil)
	if err != nil {
		log.Fatal(err)
	}
	//	v := reflect.ValueOf(r)
	//	log.Printf("type: %#v", v.Type())
	//	log.Printf("type: %#v", v.Kind())
	//	vv := reflect.ValueOf(r).Elem()

	//	for i := 0; i < vv.NumField(); i++ {
	//		vid := vv.Field(i)
	//		tid := vv.Type().Field(i)
	//		log.Printf("==================\n")
	//		//		log.Printf("%s: ", tid.Name)
	//		//		if vid.Interface() != nil {
	//		//			log.Printf("%s\n", vid.Interface())
	//		//		} else {
	//		//			log.Println("[nil]")
	//		//		}
	//		log.Printf("vid t :%T\n", vid)
	//		log.Printf("vid v :%v\n", vid)
	//		log.Printf("vid #v:%#v\n", vid)
	//		log.Printf("tid t :%T\n", tid)
	//		log.Printf("tid v :%v\n", tid)
	//		log.Printf("tid #v:%#v\n", tid)
	//	}
	//	log.Printf("%#v\n", r)
	log.Printf("Method: %s\n", r.Method)
	//	log.Printf("%v,%T,%#v\n", r.URL, r.URL, r.URL)
	log.Printf("Path: %s\n", r.URL.Path)
	log.Printf("RemoteAddr: %v\n", r.RemoteAddr)
	log.Printf("RequestURI: %v\n", r.RequestURI)
	log.Printf("Host: %v\n", r.Host)
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	for k, v := range r.Form {
		log.Printf("%s: %s\n", k, v)
	}

	for k, v := range r.Header {
		log.Printf("%s: %s\n", k, v)
	}

}
func read_config(d *ConfigData) {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(file, &d); err != nil {
		log.Fatal(err)
	}
}
func main() {
	http.HandleFunc("/start", root_handler)
	http.HandleFunc("/auth", auth_handler)
	log.Printf("start http server\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
