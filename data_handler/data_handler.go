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

type TaskForTmpl struct {
	Id int
	Project string
	Title string
	Content string
	State State
	Priority Priority
	CreatedTime string
	LastUpdateTime string
	LivedTime string
}

type TasksForTmpl struct {
	Todo []TaskForTmpl
	Doing []TaskForTmpl
	Done []TaskForTmpl
	NumberTodo int
	NumberDoing int
	NumberDone int
	NumberDoneFiltered int
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

func GetTasksForTmpl() (result TasksForTmpl) {
	tasks := readAllTasks()

	for _, task := range tasks {
		taskForHtml := ConvertTaskToTaskForTmpl(&task)

		switch taskForHtml.State {
			case Todo:
				result.Todo = append(result.Todo, taskForHtml)
				result.NumberTodo++
			case Doing:
				result.Doing = append(result.Doing, taskForHtml)
				result.NumberDoing++
			case Done:
				result.Done = append(result.Done, taskForHtml)
				result.NumberDone++
				// todo
				result.NumberDoneFiltered++
			default:
				panic("invalid state")
		}
	}

	return
}

func ConvertTaskToTaskForTmpl(task *Task) (result TaskForTmpl) {
	result.Id = task.Id
	result.Project = task.Project
	result.Title = task.Title
	result.Content = task.Content
	result.State = task.State
	result.Priority = task.Priority
	result.CreatedTime = convertTimeToString(&task.CreatedTime)
	result.LastUpdateTime = convertTimeToString(&task.LastUpdateTime)
	result.LivedTime = generatePrettyAgeForTag(task.CreatedTime)

	return
}

func generatePrettyAgeForTag(createdDate time.Time) string {
	var live_time int64 = time.Now().Sub(createdDate).Milliseconds() / 1000

    year := int64(live_time / 31536000)
    live_time = live_time - year * 31536000

    month := int64(live_time / 2592000)
    live_time = live_time - month * 2592000

    day := int64(live_time / 86400)
    live_time = live_time - day * 86400

    hour := int64(live_time / 3600)
    live_time = live_time - hour * 3600

    minute := int64(live_time / 60)

    if (year > 0) {
        return strconv.FormatInt(year, 10) + "y"
    }

    if (month > 0) {
        return strconv.FormatInt(month, 10) + "M"
    }

    if (day > 0) {
        return strconv.FormatInt(day, 10) + "d"
    }

    if (hour > 0) {
        return strconv.FormatInt(hour, 10) + "h"
    }

    return strconv.FormatInt(minute, 10) + "m"
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

func readAllTasks() []Task {
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
