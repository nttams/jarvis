const overlay = document.querySelector("#overlay");
const popup = document.querySelector("#popup");
const close_popup = document.querySelector("#close-popup");

function initHotKeys() {
    window.onload=function() {
        document.onkeyup=key_event;
    }
}

function key_event(e) {
    if (e.keyCode == 27) {
        hidePopup()
    }
}

function showPopup() {
    overlay.style.display = "block"
    popup.style.display = "block"
}

function hidePopup() {
    overlay.style.display = "none"
    popup.style.display = "none"
}

function initCreateButtion() {
    const btn = document.createElement("button")

    btn.onclick = openCreatePopup
    btn.setAttribute("id", "btn-create")
    btn.classList.add("btn")
    btn.innerHTML = "Create"

    document.querySelector("#header").appendChild(btn)
}

function initPriority() {
    const names = ["idea", "low", "med", "high"]
    const classes = ["priority-idea", "priority-low", "priority-med", "priority-high"]

    const priorities = document.querySelectorAll(".priority");

    for (let i = 0; i < priorities.length; ++i) {
        let priority_index = priorities[i].innerHTML

        priorities[i].innerHTML = names[priority_index]
        priorities[i].classList.add(classes[priority_index])

        // high priority tasks in done column don't blink
        if ((priority_index == 3) && !priorities[i].classList.contains("state-2")) {
            priorities[i].classList.add("priority-blink")
        }

    }
}

function openEditPopup(id, project, title, content, priority) {
    showPopup()
    selectPriorty(priority)

    popup.querySelector("#task-id").value = id
    popup.querySelector("#task-project").value = project
    popup.querySelector("#task-title").value =  title
    popup.querySelector("#task-content").value =  content

    popup.querySelector("#task-id-label").innerHTML = "MODIFYING TASK-" + id
    popup.querySelector(".btn-delete").disabled = false
}

function openCreatePopup() {
    showPopup()
    selectPriorty(1)

    popup.querySelector("#task-id").value = -1
    popup.querySelector("#task-project").value = ""
    popup.querySelector("#task-title").value = ""
    popup.querySelector("#task-content").value = ""

    // priority_low should be used as default, idea is only for very opaque thought
    popup.querySelector("#task-id-label").innerHTML = "CREATING NEW TASK"

    popup.querySelector("#btn-delete").disabled = true
}

function createOrUpdateTask() {
    id = document.querySelector("#task-id").value
    project = document.querySelector("#task-project").value
    title = document.querySelector("#task-title").value
    content = document.querySelector("#task-content").value
    priority = document.querySelector("#task-priority").value

    command = "create-task"
    if (id != "-1") {
        command = "update-task"
    }

    let data = {
        command: command,
        task: {
            id: parseInt(id),
            project: project,
            title: title,
            content: content,
            priority: parseInt(priority)
        }
    }

    query(data)
}


function deleteTask() {
    if (!confirm("sure?")) return;

    id = document.querySelector("#task-id").value

    let data = {
        command: "delete-task",
        task: {
            id: parseInt(id),
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
        window.location.href = "http://localhost:8080/tasks"
    })
}

function setPriorityIdea() {
    selectPriorty(0)
}

function setPriorityLow() {
    selectPriorty(1)
}


function setPriorityMed() {
    selectPriorty(2)
}


function setPriorityHigh() {
    selectPriorty(3)
}

// todo: ugly, improve this
function selectPriorty(value) {
    resetPriorty()

    document.querySelector("#task-priority").value = value

    if (value == 0 ){
        document.querySelector("#popup-priority-idea").className = "priority-blink priority-idea"
    }

    if (value == 1 ){
        document.querySelector("#popup-priority-low").className = "priority-blink priority-low"
    }

    if (value == 2 ){
        document.querySelector("#popup-priority-med").className = "priority-blink priority-med"
    }

    if (value == 3 ){
        document.querySelector("#popup-priority-high").className = "priority-blink priority-high"
    }
}

function resetPriorty() {
    document.querySelector("#popup-priority-idea").className = "priority-none"
    document.querySelector("#popup-priority-low").className = "priority-none"
    document.querySelector("#popup-priority-med").className = "priority-none"
    document.querySelector("#popup-priority-high").className = "priority-none"
}

function allowDrop(e) {
    e.preventDefault()
}

function drag(e) {
    e.dataTransfer.setData("id", e.target.innerHTML)
}

function drop(e) {
    var id = e.dataTransfer.getData("id")

    let state = -1

    if (e.target.classList.contains("dnd-state-0")) {
        state = 0
    }
    if (e.target.classList.contains("dnd-state-1")) {
        state = 1
    }
    if (e.target.classList.contains("dnd-state-2")) {
        state = 2
    }

    if (state != -1) {
        // moving to same state, do nothing
        if (document.querySelector("#" + id)
            .querySelector(".priority")
            .classList.contains("state-" + state)) {
            return
        }

        task_id = id.split("-")[1]
        let data = {
            command: "update-taks-state",
            task: {
                id: parseInt(task_id), state: state
            }
        }

        query(data)
    }
}

function init() {
    initPriority()
    initCreateButtion()
    initHotKeys()

    close_popup.onclick = hidePopup;
}

init()
