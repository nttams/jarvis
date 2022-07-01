const HOST = window.location.protocol + "//" + window.location.host + "/"
const TASK_URL = "tasks/"

const PRIORITY_IDEA = 0
const PRIORITY_LOW = 1
const PRIORITY_MED = 2
const PRIORITY_HIGH = 3
const PRIORITY_DEFAULT = PRIORITY_LOW

const STATE_TODO = 0
const STATE_DOING = 1
const STATE_DONE = 2

const OVERLAY = document.querySelector("#overlay");
const POPUP = document.querySelector("#popup");

function initHotKeys() {
    window.onload=function() {
        document.onkeyup = (e) => {
            if (e.key == "Escape") {
                hidePopup()
            }
        };
    }
}

function showPopup() {
    OVERLAY.style.display = "block"
    POPUP.style.display = "flex"
}

function hidePopup() {
    OVERLAY.style.display = "none"
    POPUP.style.display = "none"
}

function initPriority() {
    const names = ["idea", "low", "med", "high"]
    const classes = ["priority-idea", "priority-low", "priority-med", "priority-high"]

    const priorities = document.querySelectorAll(".priority");

    for (let i = 0; i < priorities.length; ++i) {
        let priority_index = priorities[i].innerHTML
        priority_index = priority_index.trim()

        priorities[i].innerHTML = names[priority_index]
        priorities[i].classList.add(classes[priority_index])
    }
}

function openEditPopup(id) {
    showPopup()
    task = getTaskById(id)

    selectPriorty(task.priority)
    POPUP.querySelector("#task-id").value = id
    POPUP.querySelector("#task-project").value = task.project
    POPUP.querySelector("#task-title").value =  task.title
    POPUP.querySelector("#task-content").value =  task.content

    POPUP.querySelector("#task-id-label").innerHTML = "editing task-" + id
}

function openCreatePopup() {
    showPopup()
    selectPriorty(PRIORITY_DEFAULT)

    POPUP.querySelector("#task-id").value = -1
    POPUP.querySelector("#task-project").value = ""
    POPUP.querySelector("#task-title").value = ""
    POPUP.querySelector("#task-content").value = ""

    // priority_low should be used as default, idea is only for very opaque thought
    POPUP.querySelector("#task-id-label").innerHTML = "creating new task"
}

function createOrUpdateTask() {
    id = parseInt(document.querySelector("#task-id").value.trim())
    project = document.querySelector("#task-project").value.trim()
    title = document.querySelector("#task-title").value.trim()
    content = document.querySelector("#task-content").value.trim()
    priority = parseInt(document.querySelector("#task-priority").value.trim())

    if (project == "" || title == "" || content == "") {
        alert("project, title, and content must not be empty")
        return
    }

    if (id != -1) {
        sendUpdateTask(id, project, title, content, priority)
    } else {
        sendCreateTask(project, title, content, priority)
    }
}

function sendCreateTask(project, title, content, priority) {
    let data = {
        command: "create-task",
        task: {
            project: project,
            title: title,
            content: content,
            priority: priority
        }
    }
    query(data)
}

function sendUpdateTask(id, project, title, content, priority) {
    let data = {
        command: "update-task",
        task: {
            id: id,
            project: project,
            title: title,
            content: content,
            priority: priority
        }
    }
    query(data)
}

function sendUpdateStateTask(id, state) {
    let data = {
        command: "update-task-state",
        task: {
            id: id,
            state: state
        }
    }
    query(data)
}

function sendDeleteTask(id) {
    let data = {
        command: "delete-task",
        task: {
            id: id,
        }
    }
    query(data)
}

function query(data) {
    fetch(HOST + TASK_URL + getLastElementInUrl(), {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data)
    }).then(res => {
        window.location.href = res.url
    })
}

// todo: ugly, improve this
function selectPriorty(value) {
    document.querySelector("#popup-priority-idea").className = "priority-none"
    document.querySelector("#popup-priority-low").className = "priority-none"
    document.querySelector("#popup-priority-med").className = "priority-none"
    document.querySelector("#popup-priority-high").className = "priority-none"

    document.querySelector("#task-priority").value = value

    switch (value) {
        case PRIORITY_IDEA:
            document.querySelector("#popup-priority-idea").className = "priority-idea"
            break;
        case PRIORITY_LOW:
            document.querySelector("#popup-priority-low").className = "priority-low"
            break;
        case PRIORITY_MED:
            document.querySelector("#popup-priority-med").className = "priority-med"
            break;
        case PRIORITY_HIGH:
            document.querySelector("#popup-priority-high").className = "priority-high"
            break;
        default:
            console.error("not recognized priority ", value)
    }
}

// since only "onDrop" can access transferData, this shared var is used for dragenver and dragleave
// todo: NOT a good solution, find something else
var current_state = -1

function drag(e, id) {
    task = getTaskById(id)
    e.dataTransfer.setData("id", id)
    e.dataTransfer.setData("project", task.project)
    e.dataTransfer.setData("title", task.title)
    e.dataTransfer.setData("state", task.state)
    current_state = task.state
}

function allowDrop(e) {
    e.preventDefault()
}

function dropOnStateColumn(e) {
    if (!e.target.classList.contains("task-column"))
        return

    let old_state = parseInt(e.dataTransfer.getData("state"))
    let state = getStateFromClassList(e.target.classList)

    if (state != old_state) {
        let id = parseInt(e.dataTransfer.getData("id"))
        sendUpdateStateTask(id, state)
        current_state = -1
    }
}

function dragEnterColumn(e) {
    if (!e.target.classList.contains("task-column"))
        return

    let state = getStateFromClassList(e.target.classList)

    if (state != current_state)
        e.target.classList.add('task-column-drag-over');
}

function getStateFromClassList(classList) {
    if (classList.contains("state-0")) return STATE_TODO
    if (classList.contains("state-1")) return STATE_DOING
    if (classList.contains("state-2")) return STATE_DONE
}

function dragLeaveColumn(e) {
    if (e.target.classList.contains("task-column"))
        e.target.classList.remove('task-column-drag-over');
}

function dropOnRecycleBin(e) {
    document.querySelector("#recycle-bin").classList.remove("recycle-bin-drag-over")
    let id = parseInt(e.dataTransfer.getData("id"))
    let project = e.dataTransfer.getData("project")
    let title = e.dataTransfer.getData("title")

    warning = `task-${id}\nproject: ${project}\ntitle: ${title}\n\ndelete?`
    if (!confirm(warning)) return;

    sendDeleteTask(id)
}

function changeProject(project) {
    project = project.split(" ")[0]

    fetch(HOST + TASK_URL + project, {
        method: "GET",
    }).then(res => {
        window.location.href = res.url
    })
}

function getLastElementInUrl() {
    url_parts = window.location.href.split("/")
    return url_parts[url_parts.length - 1]
}

function getTaskById(id) {
    for (let i = 0; i < tasks.length; ++i) {
        if (id == tasks[i].id) return tasks[i]
    }
}

function init() {
    initHotKeys()
    initPriority()
    hidePopup()

    // move tasks header to navigator
    document.querySelector("#header").appendChild(document.querySelector("#tasks-header"))
    POPUP.querySelector("#close-popup").onclick = hidePopup;

    // todo: ugly :D
    project = getLastElementInUrl()
    document.querySelector("#project-filter").value = project;
}

init()
