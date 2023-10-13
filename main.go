package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	go listenUDP(net.IPv4(127, 0, 0, 1), 12345)
	for true {
		sendUDP(net.IPv4(127, 0, 0, 1), 12345, getInput())
	}
}

func listenUDP(ip net.IP, port int) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		fmt.Println("Failed to establish udp conn.")
		return
	}
	defer conn.Close()
	for {
		var data [1024]byte
		n, _, err := conn.ReadFromUDP(data[:])
		if err != nil {
			fmt.Println("Failed to recv data from remote.")
			continue
		}

		fmt.Println("# New Msg =>", string(data[:n]))
	}
}

func getInput() string {
	fmt.Print("# Input > ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputStr := input.Text()
	return inputStr
}

func sendUDP(ip net.IP, port int, msg string) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ip,
		Port: port,
	})
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
