package main

import (
	"fmt"
	"html/template"
	"monitor"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	mm "media_manager"
	tm "task_manager"
)

var (
	templates = template.Must(template.ParseFiles(
		"tmpl/index.html",
		"tmpl/templates.html",
	))

	totalTaskRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_task_requests",
		Help: "Total Task Requests",
	})
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", 0)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	totalTaskRequests.Inc()
	tm.HandleRequest(w, r)
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	mm.HandleRequest(w, r)
}

func monitorHandler(w http.ResponseWriter, r *http.Request) {
	monitor.HandleRequest(w, r)
}

func getPort() int {
	return 8080
}

func main() {
	tm.Init()
	mm.Init()
	monitor.Init()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/media/", mediaHandler)
	http.HandleFunc("/monitor/", monitorHandler)
	http.Handle("/metrics", promhttp.Handler())

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
