function init() {
    initPopup()
    initPriority()
    initDate()
    initCreateButtion()
    initTitles()
    initHotKeys()
}

const overlay = document.querySelector("#overlay");
const popup = document.querySelector("#popup");
const close_popup = document.querySelector("#close-popup");

let names = ["Low", "Med", "High", "HOT!"]
let classes = ["priority-low", "priority-med", "priority-high", "priority-hot"]

function initHotKeys() {
    window.onload=function() {
        document.onkeyup=key_event;
    }
}

function key_event(e) {
    console.log(e.keyCode);
    if (e.keyCode == 27) {
        hidePopup()
    }
}

function initTitles() {
    numTodo = document.querySelector("#todo").childElementCount - 1
    numDoing = document.querySelector("#doing").childElementCount - 1
    numOnhold = document.querySelector("#onhold").childElementCount - 1
    numDone = document.querySelector("#done").childElementCount - 1

    document.querySelector("#todo").querySelector(".title-column").innerHTML = "Todo (" + numTodo + ")"
    document.querySelector("#doing").querySelector(".title-column").innerHTML = "Doing (" + numDoing + ")"
    document.querySelector("#onhold").querySelector(".title-column").innerHTML = "Onhold (" + numOnhold + ")"
    document.querySelector("#done").querySelector(".title-column").innerHTML = "Done (" + numDone + ")"
}

function initCreateButtion() {
    const btn = document.createElement("button")

    btn.onclick = createTask
    btn.setAttribute("id", "btn-create")
    btn.classList.add("btn")
    btn.innerHTML = "Create"

    document.querySelector("#header").appendChild(btn)
}

function initPopup() {
    close_popup.onclick = hidePopup;
}

function initPriority() {
    const priorities = document.querySelectorAll(".priority");

    for (let i = 0; i < priorities.length; ++i) {
        let priority_index = priorities[i].innerHTML

        priorities[i].classList.add(classes[priority_index])
        priorities[i].innerHTML = names[priority_index]
    }
}

function initDate() {
    const live_times = document.querySelectorAll(".live-time");

    // todo: this will not survive long-live task
    for (let i = 0; i < live_times.length; ++i) {
        let live_time = Math.abs(
            Date.now() -
            new Date(live_times[i].innerHTML)) / 1000

        day = Math.floor(live_time / 86400)

        live_time = live_time - day * 86400
        hour = Math.floor(live_time / 3600)

        if (day > 0) {
            live_times[i].innerHTML = day + "d" + hour + "h"
        } else if (hour > 0) {
            live_times[i].innerHTML = hour + "h"
        } else {
            minute = Math.floor(live_time / 60)
            live_times[i].innerHTML = minute + "m"
        }
    }
}

function editTask(id, project, title, content, state, priority) {
    showPopup()

    popup.querySelector("#task-id").value = id
    popup.querySelector("#task-project").value = project
    popup.querySelector("#task-title").value = title
    popup.querySelector("#task-content").value = content
    popup.querySelector("#task-content").value = content
    popup.querySelector("#task-id-delete").value = id

    if (state == 0) popup.querySelector("#state_todo").checked = true
    if (state == 1) popup.querySelector("#state_doing").checked = true
    if (state == 2) popup.querySelector("#state_onhold").checked = true
    if (state == 3) popup.querySelector("#state_done").checked = true

    if (priority == 0) popup.querySelector("#priority_low").checked = true
    if (priority == 1) popup.querySelector("#priority_med").checked = true
    if (priority == 2) popup.querySelector("#priority_high").checked = true
    if (priority == 3) popup.querySelector("#priority_hot").checked = true

    popup.querySelector("#task-id-label").style.display = "inline"
    popup.querySelector("#task-id-label").innerHTML = "T" + id
    popup.querySelector(".btn-delete").disabled = false
}

function createTask() {
    showPopup()

    popup.querySelector("#task-id").value = -1
    popup.querySelector("#task-title").value = ""
    popup.querySelector("#task-content").value = ""

    popup.querySelector("#state_todo").checked = true
    popup.querySelector("#priority_low").checked = true
    popup.querySelector("#task-id-label").innerHTML = "asdf"
    popup.querySelector("#task-id-label").style.display = "none"

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

init()
