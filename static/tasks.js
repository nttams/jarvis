function init() {
    const overlay = document.querySelector('#overlay');
    const popup = document.querySelector('#popup');
    const close_popup = document.querySelector('#close-popup');

    overlay.style.display = 'none'
    popup.style.display = 'none'

    close_popup.onclick = () => {
        overlay.style.display = 'none'
        popup.style.display = 'none'
    };

    const priorities = document.querySelectorAll('.priority');

    // todo: too messy?
    for (let i = 0; i < priorities.length; ++i) {
        if (priorities[i].innerHTML == 0) {
            priorities[i].innerHTML = "Low"
            priorities[i].classList.add("priority-low")
        }

        if (priorities[i].innerHTML == 1) {
            priorities[i].innerHTML = "Med"
            priorities[i].classList.add("priority-med")
        }

        if (priorities[i].innerHTML == 2) {
            priorities[i].innerHTML = "High"
            priorities[i].classList.add("priority-high")
        }

        if (priorities[i].innerHTML == 3) {
            priorities[i].innerHTML = "HOT!!"
            priorities[i].classList.add("priority-hot")
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
}

numTodo = document.querySelector("#todo").childElementCount - 1
numDoing = document.querySelector("#doing").childElementCount - 1
numOnhold = document.querySelector("#onhold").childElementCount - 1
numDone = document.querySelector("#done").childElementCount - 1

document.querySelector("#todo").querySelector(".title-column").innerHTML = "Todo (" + numTodo + ")"
document.querySelector("#doing").querySelector(".title-column").innerHTML = "Doing (" + numDoing + ")"
document.querySelector("#onhold").querySelector(".title-column").innerHTML = "Onhold (" + numOnhold + ")"
document.querySelector("#done").querySelector(".title-column").innerHTML = "Done (" + numDone + ")"

setTimeout(init, 100)