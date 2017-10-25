package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ConfigData struct {
	Client_id     string
	Client_secret string
	Redirect_uri string
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
		<body>`)
	log.SetFlags(log.Lshortfile)
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
	dumpRespBody, _ = ioutil.ReadAll(resp.Body)
//	var prof map[string]interface{}
	var prof interface{}
	if err = json.Unmarshal(dumpRespBody, &prof); err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w,"<form><textarea rows=30 cols=100>")
	assert(w,prof,0)
	fmt.Fprint(w,"</textarea></form>")
	log.Printf("=====/auth -> MailGet START=====")
//	url = "https://graph.microsoft.com/v1.0/me/mailfolders/inbox/messages?$search=public&filter=ReceiveDataTime%20ge%202017-10-01%20and%20receivedDataTime%20lt%202017-10-23&top=1"
	url = "https://graph.microsoft.com/v1.0/me/mailfolders/inbox/messages?$search=public&filter=ReceiveDataTime%20ge%202017-10-01%20and%20receivedDataTime%20lt%202017-10-23&$top=1"
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", tokens.Token_type+" "+tokens.Access_token)
	dump, _ = httputil.DumpRequestOut(req, true)
	log.Printf("request header : %s\n", dump)
	resp, err = client.Do(req)
	dumpResp, _ = httputil.DumpResponse(resp, true)
	dumpRespBody, _ = ioutil.ReadAll(resp.Body)
//	log.Printf("%s\n",dumpRespBody)
//	var dat map[string]interface{}
	var dat interface{}
	if err = json.Unmarshal(dumpRespBody, &dat); err != nil {
		log.Fatal(err)
	}
//	for k,v := range dat {
//		log.Println("==================================================")
//		log.Printf("%s = %#v\n", k,v)
//	}
fmt.Fprint(w,"<form><textarea rows=30 cols=100>")
assert(w,dat,0)
fmt.Fprint(w,"</textarea></form></body></html>")
}
func assert(w http.ResponseWriter,data interface{},idx int) {
	//インデントを生成する
	log.Printf("idx=%d",idx)
	indent_format := "@--@"
	var indent string
	for i:=0;i<idx;i++{
		indent += indent_format
	}
	
    switch data.(type) {
    case string:
        fmt.Fprintf(w,"\"%v\"\n",data.(string))
    case float64:
        fmt.Fprintf(w,"\"%v\"\n",data.(float64))
    case bool:
        fmt.Fprintf(w,"\"%v\"\n",data.(bool))
    case nil:
        fmt.Fprintf(w,"\"null\"\n")
    case []interface{}:
		fmt.Fprintf(w,"[\n")
		idx++
        for _, v := range data.([]interface{}) {
            assert(w,v,idx)
            fmt.Fprintf(w,"\"%s\"",indent)
		}
		idx--
        fmt.Fprintln(w,"]")
    case map[string]interface{}:
		fmt.Fprintf(w,"{\n")
		idx++
        for k, v := range data.(map[string]interface{}) {
            fmt.Fprintf(w,"%s\"%s\": ",indent,k)
            assert(w,v,idx)
            fmt.Fprint(w," ")
		}
		idx--
        fmt.Fprintln(w,"}")
    default:
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
	http.HandleFunc("/auth", auth_handler)
	log.Printf("start https server\n")
        err := http.ListenAndServeTLS(":8443", "ssl/service.crt","ssl/service.key",nil)
        if err != nil {
                log.Fatal(err)
        }
}
