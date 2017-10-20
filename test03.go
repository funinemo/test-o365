package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//https://oauthplay.azurewebsites.net/でトークンを取る。
//まずはそこから始めよう。
func main() {
	urls := "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=2eaea884-a699-4d42-9f1d-a6c2df1a912f&response_type=code&scope=user.read%20mail.read&state=12345&response_mode=query&redirect_uri=http%3A%2F%2Flocalhost%3A8080"

	req, _ := http.NewRequest("GET", urls, nil)
	//	log.Printf("req1 : %s\n", req)

	//	dump, _ := httputil.DumpRequestOut(req, true)
	//	log.Printf("req set header : %s\n", dump)

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	//	dumpResp, _ := httputil.DumpResponse(resp, true)
	//	log.Printf("resp : %s\n", dumpResp)
	byteArray, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("==START==")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "%s", string(byteArray))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

	//	fmt.Println(string(byteArray))
}
