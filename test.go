package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	log.Printf("==START==")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		t := template.Must(template.ParseFiles("test.html"))
		err := t.ExecuteTemplate(w, "test.html", nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("handler connect\n")
	})
	log.Printf("start http server\n")
	err := http.ListenAndServe("10.1.208.163:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("start https server\n")
	err = http.ListenAndServeTLS("10.1.208.163:8443", "ssl/service.crt","ssl/service.key",nil)
	if err != nil {
		log.Fatal(err)
	}

}
