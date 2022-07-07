const HOST = window.location.protocol + "//" + window.location.host + "/"
const TASK_URL = "monitor/192.168.1.100:9999"

function query() {
    fetch(HOST + TASK_URL)
    .then(res => {
        return res.text()
    }).then(text => {
        updateDisplay(JSON.parse(text))
    }).catch(error => console.error(error))
}
function updateDisplay(jsonData) {
    console.log(jsonData)
    document.querySelector("#cpu0").value = jsonData.CpuLoads[0]
    document.querySelector("#cpu1").value = jsonData.CpuLoads[1]
    document.querySelector("#cpu2").value = jsonData.CpuLoads[2]
    document.querySelector("#cpu3").value = jsonData.CpuLoads[3]

    // document.querySelector("#memory").value = jsonData.memory.used
    // document.querySelector("#memory").max = jsonData.memory.total

    // document.querySelector("#disk").value = jsonData.disk.used
    // document.querySelector("#disk").max = jsonData.disk.total
}

function init() {
    setInterval(query, 1000)
}

setTimeout(init, 10)
