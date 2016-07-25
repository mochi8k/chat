package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "chat.html"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = tmpl.Execute(w, r)

	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("main")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
