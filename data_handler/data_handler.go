package data_handler

import (
	"os"
	"encoding/json"
	"time"
	"strconv"
	"sort"
)

// todo: hide these fields
type Task struct {
	Id int
	Project string
	Title string
	Content string
	State State
	Priority Priority
	CreatedTime time.Time
	LastUpdateTime time.Time
}

type ByTask []Task
func (a ByTask) Len() int           { return len(a) }
func (a ByTask) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// priority (dec) --> created time (inc) --> id (inc)
func (a ByTask) Less(i, j int) bool {
	if a[i].Priority < a[j].Priority { return true }
	if a[i].Priority > a[j].Priority { return false }

	if a[i].CreatedTime.Before(a[j].CreatedTime) { return false }
	if a[i].CreatedTime.After(a[j].CreatedTime) { return true }

	if a[i].Id > a[j].Id { return true }
	return false;
}

type TaskForHtml struct {
	Id int
	Project string
	Title string
	Content string
	State State
	Priority Priority
	CreatedTime string
	LastUpdateTime string
}

type State int
const (
	Todo State = iota
	Doing
	Done
)

type Priority int
const (
	Low Priority = iota
	Med
	High
)

func ConvertTaskToTaskForHtml(task *Task) (result TaskForHtml) {
	result.Id = task.Id
	result.Project = task.Project
	result.Title = task.Title
	result.Content = task.Content
	result.State = task.State
	result.Priority = task.Priority
	result.CreatedTime = convertTimeToString(&task.CreatedTime)
	result.LastUpdateTime = convertTimeToString(&task.LastUpdateTime)

	return
}

func convertTimeToString(t *time.Time) string {
	year, month, date := t.Date()
	hour, min, _ := t.Clock()

	// minute should have 2 digits, it's prettier
	min_str := strconv.Itoa(min)
	if len(min_str) == 1 {
		min_str = "0" + min_str
	}

	return strconv.Itoa(year) + "/" +
		strconv.Itoa(int(month)) + "/" +
		strconv.Itoa(date) + " " +
		strconv.Itoa(hour) + ":" +
		min_str
}

func (s State) ToString() string {
	switch s {
	case Todo:
		return "todo"
	case Doing:
		return "doing"
	case Done:
		return "done"
	}
	return "unknown"
}

func (p Priority) ToString() string {
	switch p {
	case Low:
		return "Low"
	case Med:
		return "Med"
	case High:
		return "High"
	}
	return "unknown"
}

func saveTask(task *Task) error {
	filename := getFilename(task.Id)

	encoded, _ := json.Marshal(task)

	// todo: learn this 0600
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

func CreateNewTask(project string, title string, content string, state int, priority int) {
	id := getAFreeId()

	now := time.Now()
	task := Task {id, project, title, content, State(state), Priority(priority), now, now}
	saveTask(&task)
}

func UpdateTask(id int, project string, title string, content string, state int, priority int) {
	task := readTask(id)
	task.Project = project;
	task.Title = title;
	task.Content = content;
	task.State = State(state);
	task.Priority = Priority(priority);
	task.LastUpdateTime = time.Now();
	saveTask(&task)
}

func readTask(id int) Task {
	encoded, _ := os.ReadFile(getFilename(id))
	var task Task
	_ = json.Unmarshal(encoded, &task)
	return task
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

	sort.Sort(sort.Reverse(ByTask(result)))

	return result
}

func GetFileList(folder string) []string {
	path := "./static/" + folder
	file, _ := os.Open(path)

	files, _ := file.Readdir(0)

	result := []string{}

	for _, v := range files {
		result = append(result, path[len("."):] + "/" + v.Name())
	}
	return result
}

func getFilename(id int) string {
	return "./static/data/" + strconv.Itoa(id) + ".json"
}
