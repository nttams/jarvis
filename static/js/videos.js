function changeVideo(videoPath) {
    videoInfo = JSON.parse(videoPath)
    document.querySelector("#player").src = videoInfo.path
    document.querySelector("#player").querySelector("#track").src = videoInfo.subtitlePath
}

function init() {
    videoInfo = JSON.parse(document.querySelector("#movie-selector").value)

    document.querySelector("#player").src = videoInfo.path
    document.querySelector("#player").querySelector("#track").src = videoInfo.subtitlePath
}

init()
