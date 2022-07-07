package main

import (
	"fmt"
	"monitor"
	"net/http"
	"html/template"
	tm "task_manager"
	mm "media_manager"
)

var templates = template.Must(template.ParseFiles(
	"tmpl/index.html",
	"tmpl/templates.html",
))

func rootHandler(w http.ResponseWriter, r* http.Request) {
	templates.ExecuteTemplate(w, "index.html", 0)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tm.HandleRequest(w, r)
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	mm.HandleRequest(w, r)
}

func monitorHandler(w http.ResponseWriter, r* http.Request) {
	monitor.HandleRequest(w, r)
}

func main() {
	tm.Init()
	mm.Init()
	monitor.Init()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/media/", mediaHandler)
	http.HandleFunc("/monitor/", monitorHandler)

	http.HandleFunc("/favicon.ico", func (w http.ResponseWriter, r* http.Request) {
		http.ServeFile(w, r, "./static/img/favicon.ico")
	})
	fs_static := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))

	fs_data := http.FileServer(http.Dir("./data"))
	http.Handle("/data/", http.StripPrefix("/data/", fs_data))

	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
