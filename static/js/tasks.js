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

        priorities[i].classList.add(classes[priority_index])
        priorities[i].innerHTML = names[priority_index]
    }
}

function editTask(id, state, priority) {
    showPopup()

    let task = document.querySelector("#task-" + id)

    popup.querySelector("#task-id").value = id
    popup.querySelector("#task-project").value = task.querySelector(".task-project").innerHTML
    popup.querySelector("#task-title").value =  task.querySelector(".task-title").innerHTML
    popup.querySelector("#task-content").value =  task.querySelector(".task-content").innerHTML

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

function init() {
    initPriority()
    initCreateButtion()
    initHotKeys()

    close_popup.onclick = hidePopup;
}

init()
