function changeVideo(videoPath) {
    document.querySelector("#player").src = videoPath
}

function init() {
    document.querySelector("#player").src = document.querySelector("#movie-selector").value
}

init()
