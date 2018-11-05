package main

import (
	"./autoEntity"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseFiles("form.html", "seqForm.html"))

func newHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "form.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func seqFormHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "seqForm.html", nil)
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

func resolveSeqHandler(w http.ResponseWriter, r *http.Request) {
	tableNames := r.FormValue("tableNames")
	fmt.Printf("request tableNames are %s\n\n", tableNames)
	tableList := strings.Split(tableNames, ",")
	fmt.Println(tableList)
	code := autoEntity.GenerateSeq(tableList)
	fmt.Fprintln(w, code)
}

func main() {
	http.HandleFunc("/newTask/", newHandler)
	http.HandleFunc("/generate/", resolveHandler)
	http.HandleFunc("/generateSeq/", resolveSeqHandler)
	http.HandleFunc("/seqTask/", seqFormHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
