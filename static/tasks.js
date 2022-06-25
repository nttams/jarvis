function init() {
    initPopup()
    initPriority()
    initDate()
}

function initPopup() {
    const overlay = document.querySelector('#overlay');
    const popup = document.querySelector('#popup');
    const close_popup = document.querySelector('#close-popup');

    overlay.style.display = 'none'
    popup.style.display = 'none'

    close_popup.onclick = () => {
        overlay.style.display = 'none'
        popup.style.display = 'none'
    };
}

function initPriority() {
    const priorities = document.querySelectorAll('.priority');

    let names = ["Low", "Med", "High", "HOT!!"]
    let classes = ["priority-low", "priority-med", "priority-high", "priority-hot"]

    for (let i = 0; i < priorities.length; ++i) {
        let priority_index = priorities[i].innerHTML

        priorities[i].classList.add(classes[priority_index])
        priorities[i].innerHTML = names[priority_index]
    }
}

function initDate() {
    const live_times = document.querySelectorAll('.live-time');

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

function editTask(id, title, content, state, priority) {
    overlay.style.display = 'block'
    popup.style.display = 'block'

    document.querySelector('#popup').querySelector("#task-id").value = id
    document.querySelector('#popup').querySelector("#task-title").value = title
    document.querySelector('#popup').querySelector("#task-content").value = content
    document.querySelector('#popup').querySelector("#task-id-delete").value = id

    if (state == 0) document.querySelector('#popup').querySelector("#state_todo").checked = true
    if (state == 1) document.querySelector('#popup').querySelector("#state_doing").checked = true
    if (state == 2) document.querySelector('#popup').querySelector("#state_onhold").checked = true
    if (state == 3) document.querySelector('#popup').querySelector("#state_done").checked = true

    console.log(priority)
    if (priority == 0) document.querySelector('#popup').querySelector("#priority_low").checked = true
    if (priority == 1) document.querySelector('#popup').querySelector("#priority_med").checked = true
    if (priority == 2) document.querySelector('#popup').querySelector("#priority_high").checked = true
    if (priority == 3) document.querySelector('#popup').querySelector("#priority_hot").checked = true

    document.querySelector('#popup').querySelector("#task-id-label").innerHTML = "T-" + id
}

function createTask() {
    console.log("create new task")
    overlay.style.display = 'block'
    popup.style.display = 'block'

    document.querySelector('#popup').querySelector("#task-id").value = -1
    document.querySelector('#popup').querySelector("#task-title").value = ""
    document.querySelector('#popup').querySelector("#task-content").value = ""

    document.querySelector('#popup').querySelector("#state_todo").checked = true
    document.querySelector('#popup').querySelector("#priority_low").checked = true
    document.querySelector('#popup').querySelector("#task-id-label").innerHTML = "T-"
}

numTodo = document.querySelector("#todo").childElementCount - 1
numDoing = document.querySelector("#doing").childElementCount - 1
numOnhold = document.querySelector("#onhold").childElementCount - 1
numDone = document.querySelector("#done").childElementCount - 1

document.querySelector("#todo").querySelector(".title-column").innerHTML = "Todo " + numTodo
document.querySelector("#doing").querySelector(".title-column").innerHTML = "Doing " + numDoing
document.querySelector("#onhold").querySelector(".title-column").innerHTML = "Onhold " + numOnhold
document.querySelector("#done").querySelector(".title-column").innerHTML = "Done " + numDone

setTimeout(init, 100)
