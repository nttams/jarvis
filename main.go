package main

import (
	"fmt"
	"html/template"
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

	id, _ := strconv.Atoi(r.PostForm["id"][0])

	title := r.PostForm["title"][0]
	content := r.PostForm["content"][0]
	state, _ := strconv.Atoi(r.PostForm["state"][0])
	priority, _ := strconv.Atoi(r.PostForm["priority"][0])

	if id == -1 {
		dh.CreateNewTask(title, content, state, priority)
	} else {
		dh.UpdateTask(id, title, content, state, priority)
	}

	http.Redirect(w, r, "/tasks/", http.StatusFound)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/tasks.html")

	// todo: sort this
	tasks := dh.ReadAllTasks()

	tasksForHtml := []dh.TaskForHtml{}

	for _, task := range tasks {
		taskForHtml := dh.ConvertTaskToTaskForHtml(&task)
		tasksForHtml = append(tasksForHtml, taskForHtml)
	}

	t.Execute(w, tasksForHtml)
}

// todo: make these generic
func k50Handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/images.html")
	paths := dh.GetFileList("res/k50")

	t.Execute(w, paths)
}

func tnpdHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/images.html")
	paths := dh.GetFileList("res/tnpd")

	t.Execute(w, paths)
}

// todo:
func mediaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mediaHandler")
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/k50/", k50Handler)
	http.HandleFunc("/tnpd/", tnpdHandler)
	http.HandleFunc("/delete-task/", deleteHandler)
	http.HandleFunc("/media/", mediaHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
