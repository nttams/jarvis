const ESCAPE_KEY = 27

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
        document.onkeyup=key_event;
    }
}

function key_event(e) {
    if (e.keyCode == ESCAPE_KEY) {
        hidePopup()
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

function initCreateButtion() {
    const tasks_div = document.createElement("div")
    tasks_div.setAttribute("id", "tasks-header")

    // create button
    const btn = document.createElement("button")
    btn.onclick = openCreatePopup
    btn.setAttribute("id", "btn-create")
    btn.classList.add("btn")
    btn.innerHTML = "create"

    // create recycle bin
    const recycle_bin = document.createElement("div")
    recycle_bin.setAttribute("id", "recycle-bin")
    recycle_bin.innerHTML = "🗑"

    recycle_bin.addEventListener("drop", (event) => {
        dropOnRecycleBin(event)
    })

    recycle_bin.addEventListener("dragover", (event) => {
        allowDrop(event)
    })

    recycle_bin.addEventListener("dragenter", () => {
        recycle_bin.classList.add("recycle-bin-active")
    })

    recycle_bin.addEventListener("dragleave", () => {
        recycle_bin.classList.remove("recycle-bin-active")
    })

    tasks_div.appendChild(recycle_bin)
    tasks_div.appendChild(btn)

    document.querySelector("#header").appendChild(tasks_div)
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
    if (!confirm("sure?")) return;
    let data = {
        command: "delete-task",
        task: {
            id: id,
        }
    }
    query(data)
}

function query(data) {
    // todo: improve url, don't just copy
    fetch("http://localhost:8080/tasks/", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data)
    }).then(res => {
        // todo: go to redirect from response
        window.location.href = "http://localhost:8080/tasks"
    })
}

function setPriorityIdea() {
    selectPriorty(PRIORITY_IDEA)
}

function setPriorityLow() {
    selectPriorty(PRIORITY_LOW)
}


function setPriorityMed() {
    selectPriorty(PRIORITY_MED)
}


function setPriorityHigh() {
    selectPriorty(PRIORITY_HIGH)
}

// todo: ugly, improve this
function selectPriorty(value) {
    resetPriorty()

    document.querySelector("#task-priority").value = value

    if (value == PRIORITY_IDEA){
        document.querySelector("#popup-priority-idea").className = "priority-blink priority-idea"
    }

    if (value == PRIORITY_LOW){
        document.querySelector("#popup-priority-low").className = "priority-blink priority-low"
    }

    if (value == PRIORITY_MED){
        document.querySelector("#popup-priority-med").className = "priority-blink priority-med"
    }

    if (value == PRIORITY_HIGH){
        document.querySelector("#popup-priority-high").className = "priority-blink priority-high"
    }
}

function resetPriorty() {
    document.querySelector("#popup-priority-idea").className = "priority-none"
    document.querySelector("#popup-priority-low").className = "priority-none"
    document.querySelector("#popup-priority-med").className = "priority-none"
    document.querySelector("#popup-priority-high").className = "priority-none"
}

// drap and drop

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

function init() {
    initPriority()
    initCreateButtion()
    initHotKeys()

    POPUP.querySelector("#close-popup").onclick = hidePopup;
}

init()
