function changeWidth(width) {
    let images = document.querySelectorAll("img")

    for (let i = 0; i < images.length; i++) {
        images[i].style.width = (width - 1) +"%"
    }
}

const width = [20, 25, 50, 100]
function initOptionWidth() {
    for (let i = 0; i < width.length; ++i) {
        const option = document.createElement("option");
        option.value = width[i]
        option.innerHTML = width[i] + "%"

        document.querySelector("#width-selector").appendChild(option)
    }
    document.querySelector("#width-selector").value = width[0];
}

function init() {
    document.querySelector("#header").appendChild(document.querySelector("#images-header"))
    initOptionWidth()
    changeWidth(width[0])
}

init()
