package main

import (
	"html/template"
	"encoding/json"
	"net/http"
	"strconv"
	dh "data_handler"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getTasks(w, r)
	} else if r.Method == "POST" {
		postTask(w, r)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //todo: error handling

	id, _ := strconv.Atoi(r.PostForm["id-delete"][0])

	if id != -1 {
		dh.DeleteTask(id)
	}

	http.Redirect(w, r, "/tasks/", http.StatusFound)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //todo: error handling

	// todo: switch all to applicatin/json
	if r.Header["Content-Type"][0] == "application/x-www-form-urlencoded" {
		id, _ := strconv.Atoi(r.PostForm["id"][0])
		project := r.PostForm["project"][0]
		title := r.PostForm["title"][0]
		content := r.PostForm["content"][0]
		state, _ := strconv.Atoi(r.PostForm["state"][0])
		priority, _ := strconv.Atoi(r.PostForm["priority"][0])

		if id == -1 {
			dh.CreateNewTask(project, title, content, state, priority)
		} else {
			dh.UpdateTask(id, project, title, content, state, priority)
		}
	} else if r.Header["Content-Type"][0] == "application/json" {
		var task dh.Task
		_ = json.NewDecoder(r.Body).Decode(&task)

		dh.UpdateTaskState(task.Id, int(task.State))
	}

	http.Redirect(w, r, "/tasks/", http.StatusFound)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/html/tasks.html", "static/html/templates.html")

	t.Execute(w, dh.GetTasksForTmpl())
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
	http.HandleFunc("/delete-task/", deleteHandler)
	http.HandleFunc("/media/", mediaHandler)
	http.HandleFunc("/favicon.ico", iconHandler)

	http.ListenAndServe(":8080", nil)
}
