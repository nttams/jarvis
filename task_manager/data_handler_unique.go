/**
 * DataHandler (deprecated) uses one file for each task data
 * DataHandlerUnique uses only one file to store all task data
*/

package task_manager

import (
	"os"
	"time"
	"encoding/json"
)

const PATH = "./data/task_manager_data/all.json"

type DataHandlerUnique struct {

}

func (dh DataHandlerUnique) getAFreeId() int {
	tasks := dh.readAllTasks()
	max := -1
	for _, task := range tasks {
		if task.Id > max {
			max = task.Id
		}
	}
	return max + 1
}

func (dh *DataHandlerUnique) createTask(project string, title string, content string, priority Priority) {
	id := dh.getAFreeId()
	task := Task { id, project, title, content, Todo, Priority(priority), time.Now(), time.Now() }

	tasks := dh.readAllTasks()
	tasks = append(tasks, task)
	dh.writeAllTasks(tasks)
}

func (dh *DataHandlerUnique) updateTask(id int, project string, title string, content string, priority Priority) {
	tasks := dh.readAllTasks()
	index := findTask(tasks, id)

	// changing done task's attributes does not update lastUpdateTime
	if tasks[index].State != Done {
		tasks[index].LastUpdateTime = time.Now();
	}

	tasks[index].Project = project;
	tasks[index].Title = title;
	tasks[index].Content = content;
	tasks[index].Priority = priority;
	dh.writeAllTasks(tasks)
}

func (dh *DataHandlerUnique) changeTaskState(id int, state State) {
	tasks := dh.readAllTasks()
	index := findTask(tasks, id)

	tasks[index].State = state;
	tasks[index].LastUpdateTime = time.Now();

	dh.writeAllTasks(tasks)
}

func (dh *DataHandlerUnique) deleteTask(id int) {
	tasks := dh.readAllTasks()
	index := findTask(tasks, id)

	copy(tasks[index:], tasks[index + 1:])
	tasks = tasks[: len(tasks) - 1]

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

func (dh DataHandlerUnique) readAllTasks() []Task {
	encoded, _ := os.ReadFile(PATH)
	var tasks []Task
	json.Unmarshal(encoded, &tasks)
	return tasks
}

func (dh DataHandlerUnique) writeAllTasks(tasks []Task) {
	encoded, _ := json.MarshalIndent(tasks, "", "    ")
	os.WriteFile(PATH, encoded, 0600)
}
