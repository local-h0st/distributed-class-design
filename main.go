package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("typing anything and press enter to send.")
	// go listenUDP(net.IPv4(127, 0, 0, 1), 12345)
	go listenGroup("224.0.0.250:56789")
	for true {
		sendUDP("224.0.0.250:56789", getInput())
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
		fmt.Printf("{group msg(<%s>): %s}\n", remoteAddr, data[:n])
	}

}
