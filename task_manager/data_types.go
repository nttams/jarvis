package task_manager

import (
	"time"
)

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

// todo and doing states use this
type ByPriority []Task
func (a ByPriority) Len() int { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// isRecent --> priority (dec) --> created time (inc) --> id (inc)
func (a ByPriority) Less(i, j int) bool {
	if isRecent(a[i].LastUpdateTime) && !isRecent(a[j].LastUpdateTime) { return false }
	if !isRecent(a[i].LastUpdateTime) && isRecent(a[j].LastUpdateTime) { return true }

	if a[i].Priority < a[j].Priority { return true }
	if a[i].Priority > a[j].Priority { return false }

	if a[i].CreatedTime.Before(a[j].CreatedTime) { return false }
	if a[i].CreatedTime.After(a[j].CreatedTime) { return true }

	if a[i].Id > a[j].Id { return true }
	return false;
}

// done state uses this
type ByLastUpdate []Task
func (a ByLastUpdate) Len() int { return len(a) }
func (a ByLastUpdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// isRecent --> last update (dec) --> priority (dec) --> created time (inc) --> id (inc)
func (a ByLastUpdate) Less(i, j int) bool {
	if isRecent(a[i].LastUpdateTime) && !isRecent(a[j].LastUpdateTime) { return false }
	if !isRecent(a[i].LastUpdateTime) && isRecent(a[j].LastUpdateTime) { return true }

	if a[i].LastUpdateTime.Before(a[j].LastUpdateTime) { return true }
	if a[i].LastUpdateTime.After(a[j].LastUpdateTime) { return false }

	if a[i].Priority < a[j].Priority { return true }
	if a[i].Priority > a[j].Priority { return false }

	if a[i].CreatedTime.Before(a[j].CreatedTime) { return false }
	if a[i].CreatedTime.After(a[j].CreatedTime) { return true }

	if a[i].Id > a[j].Id { return true }
	return false;
}

type JsonRequest struct {
	Command string
	Task Task
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
	IsRecent bool
}

type ProjectInfo struct {
	Project string
	Count int
}

type ByCount []ProjectInfo
func (a ByCount) Len() int { return len(a) }
func (a ByCount) Swap(i, j int) { a[i], a[j] = a[j], a[i] } //todo: learn this swap
func (a ByCount) Less(i, j int) bool { 	return a[i].Count < a[j].Count }

type TasksForTmpl struct {
	Todo []TaskForTmpl
	Doing []TaskForTmpl
	Done []TaskForTmpl
	NumberTodo int
	NumberDoing int
	NumberDone int
	NumberDoneFiltered int
	ProjectInfos []ProjectInfo
}

type State int
const (
	Todo State = iota
	Doing
	Done
)

type Priority int
const (
	Idea Priority = iota
	Low
	Med
	High
)