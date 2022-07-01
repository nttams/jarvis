package task_manager

import (
	"time"
	"sort"
	"strconv"
	"net/http"
	"encoding/json"
	"html/template"
)

var dh DataHandlerUnique
var templates *template.Template

func Init() {
	dh = DataHandlerUnique{}
	templates = template.Must(template.ParseFiles("tmpl/tasks.html", "tmpl/templates.html"))
}

// used to parse client side data
type JsonRequest struct {
	Command string
	Task Task
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	taskGroup := r.URL.Path[len("/tasks/"):]

	if r.Method == "GET" {
		if len(taskGroup) == 0 {
			http.Redirect(w, r, "/tasks/all", http.StatusFound)
		} else {
			templates.ExecuteTemplate(w, "tasks.html", getTasksWrapper(taskGroup))
		}
	} else if r.Method == "POST" {

		var req JsonRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		switch req.Command {
			case "create-task":
				dh.createTask(req.Task.Project, req.Task.Title, req.Task.Content, req.Task.Priority)
			case "update-task":
				dh.updateTask(req.Task.Id, req.Task.Project, req.Task.Title, req.Task.Content, req.Task.Priority)
			case "update-task-state":
				dh.changeTaskState(req.Task.Id, req.Task.State)
			case "delete-task":
				dh.deleteTask(req.Task.Id)
			default:
				panic("invalid command")
		}
	}
}

func getTasksWrapper(project string) (result TasksWrapper) {
	tasks := dh.readAllTasks()
	if project == "all" {
		result = wrapTasks(tasks)
	} else {
		filteredTasks := filterProjectFromTasks(tasks, project)
		result = wrapTasks(filteredTasks)
	}

	result.ProjectInfos = collectProjectInfos(tasks)
	sort.Sort(sort.Reverse(ByCount(result.ProjectInfos)))

	return
}

func filterProjectFromTasks(tasks []Task, project string) []Task {
	var result []Task
	for _, task := range tasks {
		if task.Project == project {
			result = append(result, task)
		}
	}
	return result
}

func collectProjectInfos(tasks []Task) []ProjectInfo {
	projectInfos := []ProjectInfo{}

	for _, task := range tasks {
		found := false
		for i, projectInfo := range projectInfos {
			if task.Project == projectInfo.Name {
				projectInfos[i].Count++
				found = true
				break
			}
		}
		if !found {
			projectInfos = append(projectInfos, ProjectInfo {task.Project, 1})
		}
	}
	allProject := []ProjectInfo { ProjectInfo{"all", len(tasks)} }
	return append(allProject, projectInfos...)
}

func wrapTask(task *Task) (result TaskWrapper) {
	result.Task = *task
	result.CreatedTime = convertTimeToString(&task.CreatedTime)
	result.LastUpdateTime = convertTimeToString(&task.LastUpdateTime)
	result.LivedTime = generatePrettyAgeForTag(task.CreatedTime)
	result.IsRecent = isRecent(task.LastUpdateTime)

	return
}

// todo: you can improve this, don't loop so many times
func wrapTasks(tasks []Task) (result TasksWrapper) {
	for _, task := range tasks {
		switch task.State {
			case Idea:
				result.Idea = append(result.Idea, wrapTask(&task))
			case Todo:
				result.Todo = append(result.Todo, wrapTask(&task))
			case Doing:
				result.Doing = append(result.Doing, wrapTask(&task))
			case Done:
				result.Done = append(result.Done, wrapTask(&task))
			default:
				panic("invalid state")
		}
	}

	sort.Sort(sort.Reverse(ByPriority(result.Idea)))
	sort.Sort(sort.Reverse(ByPriority(result.Todo)))
	sort.Sort(sort.Reverse(ByPriority(result.Doing)))
	sort.Sort(sort.Reverse(ByLastUpdate(result.Done)))
	return
}

func isRecent(lastUpdateTime time.Time) bool {
	live_time := time.Now().Sub(lastUpdateTime).Milliseconds() / 1000
	// todo: config this recent time
	return live_time < 8
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
