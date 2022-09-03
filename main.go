package main

import (
	"os"
	"fmt"
	"monitor"
	"strconv"
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

func getPort() int {
	port := os.Getenv("JARVIS_PORT")
	portValue, err := strconv.Atoi(port)

	if (err != nil) {
		return 8080
	} else {
		return portValue
	}
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
	fsStatic := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))

	fsData := http.FileServer(http.Dir("./data"))
	http.Handle("/data/", http.StripPrefix("/data/", fsData))

	port := getPort()
	fmt.Println("server is listening on port:", port)
	listenAddress := ":" + strconv.Itoa(port)
	http.ListenAndServe(listenAddress, nil)
}
