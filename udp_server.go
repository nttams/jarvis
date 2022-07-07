//go:build ignore

/**
 * This program is with monitor module
 * The host needs to be monitored must deploy an instance of this program.
 * This program is a UDP server, the main web server would be the client that asks for system's
 * status, data format passing the network would be raw JSON
 * This has not been finished yet.
 */

package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"bufio"
	"strings"
	"strconv"
	"encoding/json"
)

/*
1st column : user = normal processes executing in user mode
2nd column : nice = niced processes executing in user mode
3rd column : system = processes executing in kernel mode
4th column : idle = twiddling thumbs
5th column : iowait = waiting for I/O to complete
6th column : irq = servicing interrupts
7th column : softirq = servicing softirqs
*/

type CPUInfo struct {
	Name string
	User int64
	Nice int64
	System int64
	Idle int64
	Iowait int64
	Irq int64
	Softirq int64
}

func (c *CPUInfo) sum() int64 {
	return c.User + c.Nice + c.System + c.Iowait + c.Irq + c.Softirq
}

func atoi(value string) int64 {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		panic("can not convert")
	}
	return int64(tmp)
}

func stringToCPUInfo(line string) (result CPUInfo) {
	fields := strings.Fields(line)

	result.Name = fields[0]
	result.User= atoi(fields[1])
	result.Nice= atoi(fields[2])
	result.System= atoi(fields[3])
	result.Idle= atoi(fields[4])
	result.Iowait= atoi(fields[5])
	result.Irq= atoi(fields[6])
	result.Softirq= atoi(fields[7])
	return
}

func stringToCPUInfos(lines []string) (result []CPUInfo) {
	for _, line := range lines {
		result = append(result, stringToCPUInfo(line))
	}
	return
}

func getCPUInfo() ([]CPUInfo, int64) {
	file, _ := os.Open("/proc/stat")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	now := time.Now().UnixMilli()
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			lines = append(lines, scanner.Text())
		}
	}
	return stringToCPUInfos(lines), now
}

func readCPUStatus() (result []int) {
	cpu_s, start := getCPUInfo()
	time.Sleep(300 * time.Millisecond)
	cpu_e, end := getCPUInfo()

	timeDelta := float64(end - start) //millisecond
	for i := 1; i < len(cpu_s); i++ {
		cpuDelta := float64((cpu_e[i].sum() - cpu_s[i].sum()) * 10) //millisecond
		value := int(100 * cpuDelta / timeDelta)
		result = append(result, value)
	}
	// fmt.Println()
	return
}

type SystemInfo struct {
	CpuLoads []int
}

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":9999")
	conn, _ := net.ListenUDP("udp", addr)

	for {
		buffer := make([]byte, 1024)
		_, remoteAddr, _ := conn.ReadFromUDP(buffer)


		info := SystemInfo {}
		info.CpuLoads = readCPUStatus()
		infoJson, _ := json.Marshal(info)
		fmt.Println(string(infoJson))

		conn.WriteToUDP(infoJson, remoteAddr)
	}
}
