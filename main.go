package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	tm "jarvis/task_manager"
)

var (
	templates = template.Must(template.ParseFiles(
		"tmpl/index.html",
		"tmpl/templates.html",
	))
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", 0)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tm.HandleRequest(w, r)
}

func getPort() int {
	return 8080
}

func main() {
	tm.Init()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks/", tasksHandler)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/img/favicon.ico")
	})
	fsStatic := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	port := getPort()
	fmt.Println("server is listening on port:", port)
	listenAddress := ":" + strconv.Itoa(port)
	http.ListenAndServe(listenAddress, nil)
}
