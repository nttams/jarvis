package data_handler

import (
	"os"
	"encoding/json"
	"time"
	"strconv"
)

type Task struct {
	Id int
	Title string
	Content string
	CurrentState State
	Date time.Time
}

type TaskForHtml struct {
	Id int
	Title string
	Content string
	CurrentState State
	Date string
}

type State int
const (
	Todo State = iota
	Doing
	Onhold
	Done
)

func (s State) ToString() string {
	switch s {
	case Todo:
		return "todo"
	case Doing:
		return "doing"
	case Onhold:
		return "onhold"
	case Done:
		return "done"
	}
	return "unknown"
}

func saveTask(task *Task) error {
	filename := getFilename(task.Id)

	encoded, _ := json.Marshal(task)

	return os.WriteFile(filename, encoded, 0600)
}

func DeleteTask(id int) {
	os.Remove(getFilename(id))
}

func getAFreeId() int {
	path := "./static/data"
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	max := -1
	for _, v := range files {
		filename := v.Name()[:len(v.Name())-5]
		temp, _:= strconv.Atoi(filename)
		if temp > max {
			max = temp
		}
	}

	return max + 1;
}

func CreateNewTask(title string, content string, state int) {
	id := getAFreeId()

	updateDate := time.Now()
	task := Task {id, title, content, State(state), updateDate}
	saveTask(&task)
}

func UpdateTask(id int, title string, content string, state int) {
	updateDate := time.Now()
	task := Task {id, title, content, State(state), updateDate}
	saveTask(&task)
}

func readTask(id int) (*Task, error) {
	filename := getFilename(id)
	encoded, _ := os.ReadFile(filename)

	var result Task

	_ = json.Unmarshal(encoded, &result)

	return &result, nil
}

func ReadAllTasks() []Task {
	path := "./static/data"
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	result := []Task{}

	for _, v := range files {
		encoded, _ := os.ReadFile(path + "/" + v.Name())
		var task Task
		_ = json.Unmarshal(encoded, &task)
		result = append(result, task)
	}
	return result
}

func ReadAllImagePaths(folder string) []string {
	path := "./static/res/" + folder
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	result := []string{}

	for _, v := range files {
		result = append(result, path[len("./static"):] + "/" + v.Name())
	}
	return result
}

func getFilename(id int) string {
	return "./static/data/" + strconv.Itoa(id) + ".json"
}
