package task_manager

import (
	"time"
)

type State int

const (
	Idea State = iota
	Todo
	Doing
	Done
)

type Priority int

const (
	Low Priority = iota
	Med
	High
)

type Task struct {
	Id             int       `json:"id"`
	Project        string    `json:"project"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	State          State     `json:"state"`
	Priority       Priority  `json:"priority"`
	CreatedTime    time.Time `json:"createdTime"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

type TaskWrapper struct {
	Task           Task
	CreatedTime    string
	LastUpdateTime string
	LivedTime      string
	IsRecent       bool
}

type TasksWrapper struct {
	Idea         []TaskWrapper
	Todo         []TaskWrapper
	Doing        []TaskWrapper
	Done         []TaskWrapper
	ProjectInfos []ProjectInfo
}

// used for project filter
type ProjectInfo struct {
	Name  string
	Count int
}
type ByCount []ProjectInfo

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] } //todo: learn this swap
func (a ByCount) Less(i, j int) bool { return a[i].Count < a[j].Count }

// todo and doing states use this
type ByPriority []TaskWrapper

func (a ByPriority) Len() int      { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
	// isRecent --> priority (dec) --> created time (inc) --> id (inc)
	if a[i].IsRecent && !a[j].IsRecent {
		return false
	}
	if !a[i].IsRecent && a[j].IsRecent {
		return true
	}

	if a[i].Task.Priority < a[j].Task.Priority {
		return true
	}
	if a[i].Task.Priority > a[j].Task.Priority {
		return false
	}

	if a[i].Task.CreatedTime.Before(a[j].Task.CreatedTime) {
		return false
	}
	if a[i].Task.CreatedTime.After(a[j].Task.CreatedTime) {
		return true
	}

	if a[i].Task.Id > a[j].Task.Id {
		return true
	}
	return false
}

// done state uses this
type ByLastUpdate []TaskWrapper

func (a ByLastUpdate) Len() int      { return len(a) }
func (a ByLastUpdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLastUpdate) Less(i, j int) bool {
	// isRecent --> last update (dec) --> priority (dec) --> created time (inc) --> id (inc)
	if a[i].IsRecent && !a[j].IsRecent {
		return false
	}
	if !a[i].IsRecent && a[j].IsRecent {
		return true
	}

	if a[i].Task.LastUpdateTime.Before(a[j].Task.LastUpdateTime) {
		return true
	}
	if a[i].Task.LastUpdateTime.After(a[j].Task.LastUpdateTime) {
		return false
	}

	if a[i].Task.Priority < a[j].Task.Priority {
		return true
	}
	if a[i].Task.Priority > a[j].Task.Priority {
		return false
	}

	if a[i].Task.CreatedTime.Before(a[j].Task.CreatedTime) {
		return false
	}
	if a[i].Task.CreatedTime.After(a[j].Task.CreatedTime) {
		return true
	}

	if a[i].Task.Id > a[j].Task.Id {
		return true
	}
	return false
}
