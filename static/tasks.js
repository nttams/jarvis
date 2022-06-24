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
}

function editTask(id, title, content, currentState) {
    overlay.style.display = 'block'
    popup.style.display = 'block'

    document.querySelector('#popup').querySelector("#task-id").value = id
    document.querySelector('#popup').querySelector("#task-title").value = title
    document.querySelector('#popup').querySelector("#task-content").value = content
    document.querySelector('#popup').querySelector("#task-id-delete").value = id

    if (currentState == 0) document.querySelector('#popup').querySelector("#state_todo").checked = true
    if (currentState == 1) document.querySelector('#popup').querySelector("#state_doing").checked = true
    if (currentState == 2) document.querySelector('#popup').querySelector("#state_onhold").checked = true
    if (currentState == 3) document.querySelector('#popup').querySelector("#state_done").checked = true
}

function createTask() {
    console.log("create new task")
    overlay.style.display = 'block'
    popup.style.display = 'block'

    document.querySelector('#popup').querySelector("#task-id").value = -1
    document.querySelector('#popup').querySelector("#task-title").value = ""
    document.querySelector('#popup').querySelector("#task-content").value = ""

    document.querySelector('#popup').querySelector("#state_todo").checked = true
}


numTodo = document.querySelector("#todo").childElementCount - 1
numDoing = document.querySelector("#doing").childElementCount - 1
numOnhold = document.querySelector("#onhold").childElementCount - 1
numDone = document.querySelector("#done").childElementCount - 1

document.querySelector("#todo").querySelector(".title-column").innerHTML = "Todo (" + numTodo + ")"
document.querySelector("#doing").querySelector(".title-column").innerHTML = "Todo (" + numDoing + ")"
document.querySelector("#onhold").querySelector(".title-column").innerHTML = "Todo (" + numOnhold + ")"
document.querySelector("#done").querySelector(".title-column").innerHTML = "Todo (" + numDone + ")"

setTimeout(init, 10)