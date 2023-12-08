package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/self-hosted-website", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/self-hosted-website.html"))
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/routing-with-go", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/routing-with-go.html"))
		tmpl.Execute(w, nil)
	})

	log.Println("site running on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
