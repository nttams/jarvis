package main

import (
	"fmt"
	"html/template"
	"os"
	"net/http"
	dh "data_handler"
	"strconv"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := "./wiki/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./wiki/" + title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	return &Page {Title: title, Body: body}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(w, "edit", page)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	} else {
		renderTemplate(w, "view", page)
	}
}

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
	currentState, _ := strconv.Atoi(r.PostForm["state"][0])

	if id == -1 {
		dh.CreateNewTask(title, content, currentState)
	} else {
		dh.UpdateTask(id, title, content, currentState)
	}

	http.Redirect(w, r, "/tasks/", http.StatusFound)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tasks.html")
	tasks := dh.ReadAllTasks()

	tasksForHtml := []dh.TaskForHtml{}

	for _, task := range tasks {
		year, month, date := task.Date.Date()
		hour, min, _ := task.Date.Clock()

		taskForHtml := dh.TaskForHtml{}
		taskForHtml.Id = task.Id
		taskForHtml.Title = task.Title
		taskForHtml.Content = task.Content
		taskForHtml.CurrentState = task.CurrentState
		taskForHtml.Date =
			strconv.Itoa(date) + "/" +
			strconv.Itoa(int(month)) + "/" +
			strconv.Itoa(year) + " "+
			strconv.Itoa(hour) + ":" +
			strconv.Itoa(min)

		tasksForHtml = append(tasksForHtml, taskForHtml)
	}

	t.Execute(w, tasksForHtml)
}

func k50Handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("images.html")
	paths := dh.ReadAllImagePaths("k50")

	t.Execute(w, paths)
}

func tnpdHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("images.html")
	paths := dh.ReadAllImagePaths("tnpd")

	t.Execute(w, paths)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	fmt.Println(r.FormValue("body"))
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page {Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you are looking for %s?", r.URL.Path[1:])
}

func getFreeId(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "100")
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/k50/", k50Handler)
	http.HandleFunc("/tnpd/", tnpdHandler)
	http.HandleFunc("/delete-task/", deleteHandler)

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}