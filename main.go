package main

import (
	"./autoEntity"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("form.html"))

func newHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "form.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func resolveHandler(w http.ResponseWriter, r *http.Request) {
	sql := r.FormValue("sql")
	fmt.Printf("request sql is %s\n\n", sql)
	code := autoEntity.Generate(sql)
	fmt.Fprintln(w, code)
}

func main() {
	http.HandleFunc("/newTask/", newHandler)
	http.HandleFunc("/generate/", resolveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
