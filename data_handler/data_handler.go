package data_handler

import (
	// "fmt"
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

func readTask(id int) (*Task, error) {
	filename := getFilename(id)
	encoded, _ := os.ReadFile(filename)

	var result Task

	_ = json.Unmarshal(encoded, &result)

	return &result, nil
}

func ReadAllTasks() []Task {
	path := "./data"
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

func ReadAllImagePaths() []string {
	path := "./pub/k50"
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	result := []string{}

	for _, v := range files {
		result = append(result, path[1:] + "/" + v.Name())
	}
	return result
}

func getFilename(id int) string {
	return "./data/" + strconv.Itoa(id) + ".json"
}
