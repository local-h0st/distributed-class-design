package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	var addr string
	fmt.Print("To choose mode, type 'single' or 'group':\nmode = ")
	if getInput() == "single" {
		fmt.Print("Listen single addr:port = ")
		go listenUDP(getInput())
		fmt.Print("Send to addr:port = ")
		addr = getInput()
	} else {
		fmt.Print("Group addr:port = ")
		addr = getInput()
		go listenGroup(addr)
	}

	fmt.Println("Config done.\n\nTyping anything and press Enter to send.")
	for {
		sendUDP(addr, getInput())
	}
}

func listenUDP(targetAddr string) {
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

func getInput() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputStr := input.Text()
	return inputStr
}

func sendUDP(targetAddr string, msg string) {
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

func listenGroup(targetAddr string) {
	addr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("#FromGroup(", remoteAddr, "): ", data[:n])
	}

}
