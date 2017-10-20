package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.SetFlags(log.Lmicroseconds)
		log.Println("=====START=====")
		log.SetFlags(log.Lshortfile)
		t := template.Must(template.ParseFiles("test.html"))
		err := t.ExecuteTemplate(w, "test.html", nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v\n",r)
		log.Printf("%#v\n",r)
    log.Printf("Method: %s\n", r.Method)
    log.Printf("%v,%T,%#v\n",r.URL,r.URL,r.URL)
    log.Printf("Path: %s\n",r.URL.Path)
    log.Printf("RemoteAddr: %v\n", r.RemoteAddr)
    log.Printf("RequestURI: %v\n", r.RequestURI) // ここで取れる
    log.Printf("Host: %v\n", r.Host)
    log.Printf("RawQuery: %s\n", r.URL.RawQuery) // ここでも取れる

		// だがしかし、ここで取るのが一般的だろう
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		for k,v := range r.Form {
			log.Printf("%s: %s\n", k,v)
		}
			
		for k,v := range r.Header {
			log.Printf("%s: %s\n", k,v)
		}
	})
	log.Printf("start https server\n")
	err := http.ListenAndServeTLS("10.1.208.163:8443", "ssl/service.crt","ssl/service.key",nil)
	if err != nil {
		log.Fatal(err)
	}
}
