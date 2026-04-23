package main

import (
	"fmt"
	"net"
)

func main() {

	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")} //dia chi nghe
	//mo cong ket noi bang listenUDP
	connect, _ := net.ListenUDP("udp", &addr) //tra ve UDPConn va error

	defer connect.Close() //phong truong hop bi loi thi van co the out

	fmt.Println("DNS Server is listening on port 53!!!")

	for {
		//chuan bi nhan du lieu
		buf := make([]byte, 512)
		//ReadFromUDP tra ve int(number of byte), Addr va error
		n, address, err := connect.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//if it not say anything try "nslookup google.com 127.0.0.1" in another powershell
		fmt.Printf("Received %d bytes from %s: %x\n", n, address, buf[:n])
	}

}
