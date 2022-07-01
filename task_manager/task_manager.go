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

func Init() {
	dh = DataHandlerUnique{}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	taskGroup := r.URL.Path[len("/tasks/"):]

	if r.Method == "GET" {
		if len(taskGroup) == 0 {
			http.Redirect(w, r, "/tasks/all", http.StatusFound)
		} else {
			// todo: ineffective, read once please
			// currently, I don't do that because of testcase
			// it fails if it cannot find the html file
			templates := template.Must(template.ParseFiles("tmpl/tasks.html", "tmpl/templates.html"))
			templates.ExecuteTemplate(w, "tasks.html", getAllTasksForTmpl(taskGroup))
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

func getAllTasksForTmpl(project string) (result TasksForTmpl) {
	tasks := dh.readAllTasks()
	if project != "all" {
		filteredTasks := filterProjectFromTasks(tasks, project)
		result = convertTasksToTasksForTmpl(filteredTasks)
	} else {
		result = convertTasksToTasksForTmpl(tasks)
	}

	result.ProjectInfos = collectProjectInfos(tasks)

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

func convertTasksToTasksForTmpl(tasks []Task) (result TasksForTmpl) {
	todo, doing, done := divideTasksIntoGroups(tasks)

	sort.Sort(sort.Reverse(ByPriority(todo)))
	sort.Sort(sort.Reverse(ByPriority(doing)))
	sort.Sort(sort.Reverse(ByLastUpdate(done)))

	for _, task := range todo {
		result.Todo = append(result.Todo, convertTaskToTaskForTmpl(&task))
	}

	for _, task := range doing {
		result.Doing = append(result.Doing, convertTaskToTaskForTmpl(&task))
	}

	for _, task := range done {
		result.Done = append(result.Done, convertTaskToTaskForTmpl(&task))
	}

	result.NumberTodo = len(todo)
	result.NumberDoing = len(doing)
	result.NumberDone = len(done)
	result.NumberDoneFiltered = len(done)

	return
}

func getDistinctProjects(tasks []Task) (result []string) {
	for _, task := range tasks {
		if !isIn(result, task.Project) {
			result = append(result, task.Project)
		}
	}
	return
}

func countProject(tasks []Task, project string) int {
	result := 0;
	for _, task := range tasks {
		if (task.Project == project) {
			result++
		}
	}
	return result
}

// todo: ugly
func collectProjectInfos(tasks []Task) (result []ProjectInfo) {
	distinctProjects := getDistinctProjects(tasks)

	result = append(result, ProjectInfo {"all", 0})

	for _, project := range distinctProjects {
		result = append(result, ProjectInfo {project, 0})
	}

	for i, _ := range result {
		result[i].Count = countProject(tasks, result[i].Project)
	}

	// todo: test if this copies or points. i'm lazy now
	// for _, projectInfo := range result {

	result[0].Count = len(tasks)

	sort.Sort(sort.Reverse(ByCount(result)))

	return
}


func isIn(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true;
		}
	}
	return false;
}

func divideTasksIntoGroups(tasks []Task) (todo []Task, doing []Task, done[]Task){
	for _, task := range tasks {
		switch task.State {
			case Todo:
				todo = append(todo, task)
			case Doing:
				doing = append(doing, task)
			case Done:
				done = append(done, task)
			default:
				panic("invalid state")
		}
	}
	return
}

func convertTaskToTaskForTmpl(task *Task) (result TaskForTmpl) {
	result.Id = task.Id
	result.Project = task.Project
	result.Title = task.Title
	result.Content = task.Content
	result.State = task.State
	result.Priority = task.Priority
	result.CreatedTime = convertTimeToString(&task.CreatedTime)
	result.LastUpdateTime = convertTimeToString(&task.LastUpdateTime)
	result.LivedTime = generatePrettyAgeForTag(task.CreatedTime)

	result.IsRecent = isRecent(task.LastUpdateTime)

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
