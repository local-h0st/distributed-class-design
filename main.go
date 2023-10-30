package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func getInput() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputStr := input.Text()
	return inputStr
}

type HOMEWORK_INFO struct {
	Name  string
	ID    string
	Seq   string
	Grade string
	Tag   string
}

func main() {
	fmt.Print("To choose mode, type 'single' or 'group':\nmode = ")
	if getInput() == "single" {
		fmt.Print("Listen single addr:port = ")
		go listenSingle(getInput())
		fmt.Print("Send to addr:port = ")
		raddr := getInput()
		fmt.Println("Config done. Typing anything and press Enter to send.")
		for {
			sendMsg(raddr, getInput())
		}
	} else {
		fmt.Print("Group addr:port = ")
		gaddr := getInput()
		go listenGroup(gaddr)
		fmt.Println("Config done. Typing anything and press Enter to send.")
		for {
			sendMsg(gaddr, getInput())
		}
	}

}

func sendMsg(targetAddr string, msg string) {
	addr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		fmt.Println(err)
	}
	socket, err := net.DialUDP("udp", nil, addr)
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

func listenSingle(targetAddr string) {
	addr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Failed to establish udp conn.")
		return
	}
	defer conn.Close()
	for {
		var data [1024]byte
		n, remoteAddr, err := conn.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println("Failed to recv data from remote.")
			continue
		}

		fmt.Println("#FromSingle(", remoteAddr, "): ", string(data[:n]))
	}
}

func listenGroup(groupAddr string) {
	gaddr, _ := net.ResolveUDPAddr("udp4", groupAddr)
	conn, _ := net.ListenMulticastUDP("udp", nil, gaddr)
	conn.SetReadBuffer(1024)
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf(string(remoteAddr.IP) + ":" + string(remoteAddr.Port) + "-> " + string(data[:n]) + "\n")
	}
}
