package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	http.Handle("/static/",
    http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("./templates/index.html"))
    tmpl.Execute(w, nil)
  })

  log.Println("Site running on port 3000...")
  log.Fatal(http.ListenAndServe(":3000", nil))
}
