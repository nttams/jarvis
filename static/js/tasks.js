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

function initCreateButtion() {
    const btn = document.createElement("button")

    btn.onclick = createTask
    btn.setAttribute("id", "btn-create")
    btn.classList.add("btn")
    btn.innerHTML = "Create"

    document.querySelector("#header").appendChild(btn)
}

function initPriority() {
    const names = ["low", "med", "high"]
    const classes = ["priority-low", "priority-med", "priority-high"]

    const priorities = document.querySelectorAll(".priority");

    for (let i = 0; i < priorities.length; ++i) {
        let priority_index = priorities[i].innerHTML

        priorities[i].innerHTML = names[priority_index]
        priorities[i].classList.add(classes[priority_index])

        // high priority tasks in done column don't blink
        if ((priority_index == 2) && !priorities[i].classList.contains("state-2")) {
            priorities[i].classList.add("priority-blink")
        }

    }
}

function editTask(id, project, title, content, state, priority) {
    showPopup()

    popup.querySelector("#task-id").value = id
    popup.querySelector("#task-project").value = project
    popup.querySelector("#task-title").value =  title
    popup.querySelector("#task-content").value =  content

    popup.querySelector("#task-id-delete").value = id

    if (state == 0) popup.querySelector("#state_todo").checked = true
    if (state == 1) popup.querySelector("#state_doing").checked = true
    if (state == 2) popup.querySelector("#state_done").checked = true

    if (priority == 0) popup.querySelector("#priority_low").checked = true
    if (priority == 1) popup.querySelector("#priority_med").checked = true
    if (priority == 2) popup.querySelector("#priority_high").checked = true

    popup.querySelector("#task-id-label").innerHTML = "MODIFYING TASK-" + id
    popup.querySelector(".btn-delete").disabled = false
}

function createTask() {
    showPopup()

    popup.querySelector("#task-id").value = -1
    popup.querySelector("#task-project").value = ""
    popup.querySelector("#task-title").value = ""
    popup.querySelector("#task-content").value = ""

    popup.querySelector("#state_todo").checked = true
    popup.querySelector("#priority_low").checked = true
    popup.querySelector("#task-id-label").innerHTML = "CREATING NEW TASK"

    popup.querySelector(".btn-delete").disabled = true
}

function showPopup() {
    overlay.style.display = "block"
    popup.style.display = "block"
}

function hidePopup() {
    overlay.style.display = "none"
    popup.style.display = "none"
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
        task_id = id.split("-")[1]
        let data = {id: parseInt(task_id), state: state}

        if (document.querySelector("#" + id)
            .querySelector(".priority")
            .classList.contains("state-" + state)) {
            return
        }

        // todo: improve url, don't just copy
        fetch("http://localhost:8080/tasks/", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(data)
        }).then(res => {
            window.location.href = "http://localhost:8080/tasks"
        })
    }

}

function init() {
    initPriority()
    initCreateButtion()
    initHotKeys()

    close_popup.onclick = hidePopup;
}

init()
