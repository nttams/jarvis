package task_manager

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const PATH = "./data/task_manager_data/tasks.json"

type DataHandler struct {
}

func (dh DataHandler) getAFreeId() int {
	tasks := dh.readAllTasks()
	max := -1
	for _, task := range tasks {
		if task.Id > max {
			max = task.Id
		}
	}
	return max + 1
}

func (dh *DataHandler) createTask(project string, title string, content string, priority Priority) {
	id := dh.getAFreeId()
	task := Task{id, project, title, content, Idea, Priority(priority), time.Now(), time.Now()}

	tasks := dh.readAllTasks()
	tasks = append(tasks, task)
	dh.writeAllTasks(tasks)
}

func (dh *DataHandler) updateTask(id int, project string, title string, content string, priority Priority) {
	tasks := dh.readAllTasks()
	index := findTask(tasks, id)

	// changing done task's attributes does not update lastUpdateTime
	if tasks[index].State != Done {
		tasks[index].LastUpdateTime = time.Now()
	}

	tasks[index].Project = project
	tasks[index].Title = title
	tasks[index].Content = content
	tasks[index].Priority = priority
	dh.writeAllTasks(tasks)
}

func (dh *DataHandler) changeTaskState(id int, state State) {
	tasks := dh.readAllTasks()
	index := findTask(tasks, id)

	tasks[index].State = state
	tasks[index].LastUpdateTime = time.Now()

	dh.writeAllTasks(tasks)
}

func (dh *DataHandler) deleteTask(id int) {
	tasks := dh.readAllTasks()
	index := findTask(tasks, id)

	copy(tasks[index:], tasks[index+1:])
	tasks = tasks[:len(tasks)-1]

	dh.writeAllTasks(tasks)
}

func findTask(tasks []Task, id int) int {
	for i, task := range tasks {
		if task.Id == id {
			return i
		}
	}
	panic("invalid index")
}

func (dh DataHandler) readAllTasks() []Task {
	encoded, err := os.ReadFile(PATH)
	if err != nil {
		log.Fatalf("failed to read file: %s, error: %s", PATH, err)
	}
	var tasks []Task
	json.Unmarshal(encoded, &tasks)
	return tasks
}

func (dh DataHandler) writeAllTasks(tasks []Task) {
	encoded, _ := json.MarshalIndent(tasks, "", "    ")
	os.WriteFile(PATH, encoded, 0600)
}
