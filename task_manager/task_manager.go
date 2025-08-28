package task_manager

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var dh DataHandler
var templates *template.Template

func Init() {
	dh = DataHandler{}
	templates = template.Must(template.ParseFiles("tmpl/tasks.html", "tmpl/templates.html"))
}

// used to parse client side data
type JsonRequest struct {
	Command string
	Task    Task
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL, "from", r.RemoteAddr)

	taskGroup := r.URL.Path[len("/tasks/"):]

	switch r.Method {
	case http.MethodGet:
		if len(taskGroup) == 0 {
			http.Redirect(w, r, "/tasks/all", http.StatusFound)
		} else {
			templates.ExecuteTemplate(w, "tasks.html", getTasksWrapper(taskGroup))
		}
	case http.MethodPost:
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
			projectInfos = append(projectInfos, ProjectInfo{task.Project, 1})
		}
	}
	allProject := []ProjectInfo{{"all", len(tasks)}}
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
	liveTime := time.Since(lastUpdateTime).Milliseconds() / 1000
	// todo: config this recent time
	return liveTime < 8
}

func generatePrettyAgeForTag(createdDate time.Time) string {
	var liveTime int64 = time.Since(createdDate).Milliseconds() / 1000

	year := int64(liveTime / 31536000)
	liveTime = liveTime - year*31536000

	month := int64(liveTime / 2592000)
	liveTime = liveTime - month*2592000

	day := int64(liveTime / 86400)
	liveTime = liveTime - day*86400

	hour := int64(liveTime / 3600)
	liveTime = liveTime - hour*3600

	minute := int64(liveTime / 60)

	if year > 0 {
		return strconv.FormatInt(year, 10) + "y"
	}

	if month > 0 {
		return strconv.FormatInt(month, 10) + "M"
	}

	if day > 0 {
		return strconv.FormatInt(day, 10) + "d"
	}

	if hour > 0 {
		return strconv.FormatInt(hour, 10) + "h"
	}

	return strconv.FormatInt(minute, 10) + "m"
}

func convertTimeToString(t *time.Time) string {
	year, month, date := t.Date()
	hour, min, _ := t.Clock()

	minStr := strconv.Itoa(min)
	if len(minStr) == 1 {
		minStr = "0" + minStr
	}

	return strconv.Itoa(year) + "/" +
		strconv.Itoa(int(month)) + "/" +
		strconv.Itoa(date) + " " +
		strconv.Itoa(hour) + ":" +
		minStr
}
