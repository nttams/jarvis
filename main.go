package main

import (
	"html/template"
	"encoding/json"
	"net/http"
	dh "data_handler"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// todo: ignore all details url?
		getAllTasks(w)
	} else if r.Method == "POST" {

		var req dh.JsonRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		switch req.Command {
			case "create-task":
				dh.CreateTask(req.Task.Project, req.Task.Title, req.Task.Content, req.Task.Priority)
			case "update-task":
				dh.UpdateTask(req.Task.Id, req.Task.Project, req.Task.Title, req.Task.Content, req.Task.Priority)
			case "update-task-state":
				dh.UpdateTaskState(req.Task.Id, req.Task.State)
			case "delete-task":
				dh.DeleteTask(req.Task.Id)
			default:
				panic("invalid command")
		}
	}
}

func getAllTasks(w http.ResponseWriter) {
	t, _ := template.ParseFiles("static/html/tasks.html", "static/html/templates.html")

	t.Execute(w, dh.GetAllTasksForTmpl())
}

// todo: make these generic
func k50Handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/images.html", "static/html/templates.html")
	paths := dh.GetFileList("res/k50")

	t.Execute(w, paths)
}

func tnpdHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/images.html", "static/html/templates.html")
	paths := dh.GetFileList("res/tnpd")

	t.Execute(w, paths)
}

// todo:
func mediaHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/media.html", "static/html/templates.html")
	t.Execute(w, 0)
}

func iconHandler(w http.ResponseWriter, r* http.Request) {
	http.ServeFile(w, r, "./static/img/favicon.ico")
}

func rootHandler(w http.ResponseWriter, r* http.Request) {
	t, _ := template.ParseFiles("static/html/index.html", "static/html/templates.html")
	t.Execute(w, 0)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/k50/", k50Handler)
	http.HandleFunc("/tnpd/", tnpdHandler)
	http.HandleFunc("/media/", mediaHandler)
	http.HandleFunc("/favicon.ico", iconHandler)

	http.ListenAndServe(":8080", nil)
}
