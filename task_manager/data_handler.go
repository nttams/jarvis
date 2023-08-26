// deprecated, use DataHandlerUnique instead

package task_manager

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

const DATA_PATH = "./data/task_manager_data/data/"

// DataHandler should be stateful for caching stuff
// But not for now :D
type DataHandler struct {
}

func getAFreeId() int {
	file, _ := os.Open(DATA_PATH)

	files, _ := file.Readdir(0)

	max := -1
	for _, v := range files {
		filename := v.Name()[:len(v.Name())-5]
		temp, _ := strconv.Atoi(filename)
		if temp > max {
			max = temp
		}
	}
	return max + 1
}

func (dh *DataHandler) createTask(project string, title string, content string, priority Priority) {
	id := getAFreeId()

	now := time.Now()
	task := Task{id, project, title, content, Idea, Priority(priority), now, now}
	writeTask(&task)
}

func (dh *DataHandler) updateTask(id int, project string, title string, content string, priority Priority) {
	task := readTask(id)

	// changing attribute in done task does not update lastUpdateTime
	if task.State != Done {
		task.LastUpdateTime = time.Now()
	}

	task.Project = project
	task.Title = title
	task.Content = content
	task.Priority = priority
	writeTask(&task)
}

func (dh *DataHandler) changeTaskState(id int, state State) {
	task := readTask(id)

	task.State = state
	task.LastUpdateTime = time.Now()
	writeTask(&task)
}

func (dh *DataHandler) readAllTasks() []Task {
	file, _ := os.Open(DATA_PATH)

	files, _ := file.Readdir(0)

	result := []Task{}

	for _, v := range files {
		encoded, _ := os.ReadFile(DATA_PATH + v.Name())
		var task Task
		_ = json.Unmarshal(encoded, &task)
		result = append(result, task)
	}

	return result
}

func (dh *DataHandler) deleteTask(id int) {
	os.Remove(getFilePath(id))
}

func getFilePath(id int) string {
	return DATA_PATH + strconv.Itoa(id) + ".json"
}

func writeTask(task *Task) error {
	filename := getFilePath(task.Id)

	encoded, _ := json.Marshal(task)

	// todo: learn this 0600
	return os.WriteFile(filename, encoded, 0600)
}

func readTask(id int) Task {
	encoded, _ := os.ReadFile(getFilePath(id))
	var task Task
	json.Unmarshal(encoded, &task)
	return task
}
