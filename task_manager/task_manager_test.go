package task_manager

import (
	"time"
	"testing"
)

func TestGetDistinctProject(t *testing.T) {
	task0 := Task {0, "project0", "title0", "content0", Idea, Low, time.Now(), time.Now()}
	task1 := Task {1, "project0", "title0", "content0", Todo, Low, time.Now(), time.Now()}

	task2 := Task {2, "project2", "title0", "content0", Doing, Low, time.Now(), time.Now()}

	tasks := []Task {}
	tasks = append(tasks, task0)
	tasks = append(tasks, task1)
	tasks = append(tasks, task2)

	projects := getDistinctProject(tasks)

	if len(projects) != 2 || projects[0] != "project0" || projects[1] != "project2" {
		t.Fatal("failed")
	}
}
