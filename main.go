package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type HOMEWORK_INFO struct {
	Name  string
	ID    string
	Seq   string
	Grade string
	Tag   string
}

var localAddr string
var groupAddr string
var HomeworkRecords []HOMEWORK_INFO

func main() {
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
		fmt.Println("Send successfully, myself received.")
	} else {
		h := HOMEWORK_INFO{}
		err := json.Unmarshal([]byte(jsonstr), &h)
		if err != nil {
			fmt.Println("Msg json received but unmarshal failed: " + jsonstr)
			return
		}
		HomeworkRecords = append(HomeworkRecords, h)
		fmt.Println("Record received from " + raddr.String() + ", and it has been stored to slice []HomeworkRecords:")
		fmt.Println(HomeworkRecords)
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
