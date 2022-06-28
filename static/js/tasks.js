const HOST = "http://localhost:8080/"
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
    POPUP.style.display = "block"
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

        // high priority tasks in done column don't blink
        if ((priority_index == PRIORITY_HIGH) && !priorities[i].classList.contains("state-2")) {
            priorities[i].classList.add("priority-blink")
        }

    }
}

function openEditPopup(id, project, title, content, priority) {
    showPopup()
    selectPriorty(priority)

    POPUP.querySelector("#task-id").value = id
    POPUP.querySelector("#task-project").value = project
    POPUP.querySelector("#task-title").value =  title
    POPUP.querySelector("#task-content").value =  content

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
    id = parseInt(document.querySelector("#task-id").value)
    project = document.querySelector("#task-project").value
    title = document.querySelector("#task-title").value
    content = document.querySelector("#task-content").value
    priority = parseInt(document.querySelector("#task-priority").value)

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
    if (!confirm("deleting task-" + id + ", sure?")) return;
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
            document.querySelector("#popup-priority-idea").className = "priority-blink priority-idea"
            break;
        case PRIORITY_LOW:
            document.querySelector("#popup-priority-low").className = "priority-blink priority-low"
            break;
        case PRIORITY_MED:
            document.querySelector("#popup-priority-med").className = "priority-blink priority-med"
            break;
        case PRIORITY_HIGH:
            document.querySelector("#popup-priority-high").className = "priority-blink priority-high"
            break;
        default:
            console.error("not recognized priorty")
    }
}

function drag(e, id) {
    e.dataTransfer.setData("id", id)
}

function allowDrop(e) {
    e.preventDefault()
}

function dropOnStateColumn(e) {
    var id = parseInt(e.dataTransfer.getData("id"))

    let state = -1

    if (e.target.classList.contains("dnd-state-0")) {
        state = STATE_TODO
    }
    if (e.target.classList.contains("dnd-state-1")) {
        state = STATE_DOING
    }
    if (e.target.classList.contains("dnd-state-2")) {
        state = STATE_DONE
    }

    if (state == -1) return

    // moving to same state, do nothing
    if (document.querySelector("#T-" + id)
        .querySelector(".priority")
        .classList.contains("state-" + state)) {
        return
    }

    sendUpdateStateTask(id, state)
}

function dropOnRecycleBin(e) {
    document.querySelector("#recycle-bin").classList.remove("recycle-bin-active")
    var id = parseInt(e.dataTransfer.getData("id"))
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

function init() {
    initHotKeys()
    initPriority()

    // move tasks header to navigator
    document.querySelector("#header").appendChild(document.querySelector("#tasks-header"))
    POPUP.querySelector("#close-popup").onclick = hidePopup;

    // todo: ugly :D
    project = getLastElementInUrl()
    document.querySelector("#project-filter").value = project;
}

init()
