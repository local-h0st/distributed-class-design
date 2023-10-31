package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type HOMEWORK_INFO struct {
	// 发送、接受消息的格式
	Name  string
	ID    string
	Seq   string
	Grade string
	Tag   string
}

type STATIC_INFO struct {
	// 每次作业统计数据格式
	Total       int
	InTimeCount int
	Sum         float64
	Highest     float64
}

var localAddr string
var groupAddr string
var HomeworkRecords []HOMEWORK_INFO
var statics map[string]*STATIC_INFO

func main() {
	statics = make(map[string]*STATIC_INFO)
	fmt.Print("Set Multi-cast group addr:port =")
	groupAddr = getInput()
	fmt.Print("Set local port = ")
	localAddr = getInterfaceIP("eth0") + ":" + getInput()
	go listenGroup(groupAddr, handleInfo)
	for {
		sendStr(groupAddr, generateHomeworkInfo())
	}
}

func handleInfo(jsonstr string, raddr *net.UDPAddr) {
	if raddr.String() == localAddr {
		fmt.Println("Send successfully(myself received).")
		//return
	}

	h := HOMEWORK_INFO{}
	err := json.Unmarshal([]byte(jsonstr), &h)
	if err != nil {
		fmt.Println("Msg json received but unmarshal failed: " + jsonstr)
		return
	}
	HomeworkRecords = append(HomeworkRecords, h)
	fmt.Println("Record received from " + raddr.String())
	fmt.Println("[]HomeworkRecords =", HomeworkRecords)

	// update the statics
	_, keyExists := statics[h.Seq]
	if !keyExists {
		statics[h.Seq] = &STATIC_INFO{
			Total:       0,
			InTimeCount: 0,
			Sum:         0,
			Highest:     0,
		}
	}
	statics[h.Seq].Total++
	grade, _ := strconv.ParseFloat(h.Grade, 64)
	statics[h.Seq].Sum += grade
	if grade > statics[h.Seq].Highest {
		statics[h.Seq].Highest = grade
	}
	if h.Tag == "Yes" {
		statics[h.Seq].InTimeCount++
	}

	// print the statics
	for k, static := range statics {
		fmt.Println("Homework", k, "->", static.InTimeCount, "/", static.Total)
		fmt.Println("Average:", static.Sum/float64(static.Total), "\tHighest:", static.Highest)
	}
}

func getInterfaceIP(name string) string {
	ifi, _ := net.InterfaceByName(name)
	addrs, _ := ifi.Addrs()
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func getInput() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputStr := input.Text()
	return inputStr
}

func generateHomeworkInfo() string {
	fmt.Println("Enter Name, ID, Seq, Grade, Tag in sequence.")
	h := HOMEWORK_INFO{
		Name:  getInput(),
		ID:    getInput(),
		Seq:   getInput(),
		Grade: getInput(),
		Tag:   getInput(),
	}
	s, _ := json.Marshal(h)
	return string(s)
}

func sendStr(targetAddr string, msg string) {
	raddr, _ := net.ResolveUDPAddr("udp", targetAddr)
	laddr, _ := net.ResolveUDPAddr("udp", localAddr)
	socket, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Println("Failed to connect target addr.")
		return
	}
	defer socket.Close()
	sendData := []byte(msg)
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("Failed to send msg.")
		return
	}
}
func listenGroup(groupAddr string, handleFunc func(data string, raddr *net.UDPAddr)) {
	gaddr, _ := net.ResolveUDPAddr("udp4", groupAddr)
	conn, _ := net.ListenMulticastUDP("udp", nil, gaddr)
	conn.SetReadBuffer(1024)
	for {
		data := make([]byte, 1024)
		n, raddr, _ := conn.ReadFromUDP(data)
		handleFunc(string(data[:n]), raddr)
	}
}
