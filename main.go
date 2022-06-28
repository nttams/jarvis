package main

import (
	"html/template"
	"encoding/json"
	"net/http"
	dh "data_handler"
)

var tasks_templates = template.Must(template.ParseFiles(
		"tmpl/index.html",
		"tmpl/tasks.html",
		"tmpl/media.html",
		"tmpl/images.html",
		"tmpl/templates.html",
	))

func rootHandler(w http.ResponseWriter, r* http.Request) {
	tasks_templates.ExecuteTemplate(w, "index.html", 0)
}

func iconHandler(w http.ResponseWriter, r* http.Request) {
	http.ServeFile(w, r, "./static/img/favicon.ico")
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path[len("/tasks/"):]

	if r.Method == "GET" {
		if len(parts) == 0 {
			http.Redirect(w, r, "/tasks/all", http.StatusFound)
		} else {
			tasks_templates.ExecuteTemplate(w, "tasks.html", dh.GetAllTasksForTmpl(parts))
		}
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

// todo: make these generic
func k50Handler(w http.ResponseWriter, r *http.Request) {
	tasks_templates.ExecuteTemplate(w, "images.html", dh.GetFileList("res/k50"))
}

func tnpdHandler(w http.ResponseWriter, r *http.Request) {
	tasks_templates.ExecuteTemplate(w, "images.html", dh.GetFileList("res/tnpd"))
}

// todo:
func mediaHandler(w http.ResponseWriter, r *http.Request) {
	tasks_templates.ExecuteTemplate(w, "media.html", 0)
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/favicon.ico", iconHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/media/", mediaHandler)
	http.HandleFunc("/k50/", k50Handler)
	http.HandleFunc("/tnpd/", tnpdHandler)

	http.ListenAndServe(":8080", nil)
}
