package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lp := filepath.Join("templates", "index.html")
		fp := filepath.Join("templates", "articles.html")
		tmpl, _ := template.ParseFiles(lp, fp)
		tmpl.ExecuteTemplate(w, "index", nil)
	})

	mux.HandleFunc("/{article}", func(w http.ResponseWriter, r *http.Request) {
		article := r.PathValue("article")
		lp := filepath.Join("templates", "article-layout.html")
		fp := filepath.Join("templates", article+".html")
		tmpl, _ := template.ParseFiles(lp, fp)
		tmpl.ExecuteTemplate(w, "article-layout", nil)
	})

	log.Println("site running on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
