package routes

import (
	"log"
	"net/http"
	"runtime/debug"
	"text/template"
)

func PageHomeHandler(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var err error

	if r.URL.Path != "/" && r.URL.Path != "/favicon.ico" {
		log.Println("PageHomeHandler::Incorrect error path ", r.URL.Path)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
