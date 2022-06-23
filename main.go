package main

import (
	"fmt"
	"html/template"
	"os"
	"net/http"
	"data_handler"
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
	t, _ := template.ParseFiles("tasks.html")
	tasks := data_handler.ReadAllTasks()

	t.Execute(w, tasks)
}

func k50Handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("images.html")
	paths := data_handler.ReadAllImagePaths()

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

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/k50/", k50Handler)

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}