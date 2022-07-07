package monitor

import (
	"fmt"
	"net"
	"bufio"
	"net/http"
	"html/template"
)

var templates *template.Template

func Init() {
	templates = template.Must(template.ParseFiles("tmpl/monitor.html", "tmpl/templates.html"))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	requestHost := r.URL.Path[len("/monitor/"):]

	if r.Method == "GET" {
		if len(requestHost) == 0 {
			templates.ExecuteTemplate(w, "monitor.html", 0)
		} else {
			remoteHost, _ := net.ResolveUDPAddr("udp", requestHost)
			conn, _ := net.DialUDP("udp", &net.UDPAddr {}, remoteHost)
			conn.Write([]byte("anything"))
			buffer:= make([]byte, 1024)
			count, _ := bufio.NewReader(conn).Read(buffer)
			fmt.Println(buffer[:count])
			fmt.Println()
			fmt.Fprintf(w, string(buffer[:count]))
		}
	} else {
		panic("does not support POST in monitor")
	}
}
