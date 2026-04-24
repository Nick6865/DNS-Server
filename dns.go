package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const banner = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣶⣶⣶⣶⣶⣶⣦⣄⡀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⡀⠀⠀⠀⠀⠀⢀⣈⣿⣶⣿⣿⣿⡿⢿⣿⡿⠛⠙⠻⣿⣿⡇⠀⠀
⠀⠀⠀⠀⠀⣀⣤⣤⣶⣶⣆⣀⣤⣄⡀⠀⠀⠀⢀⣤⣶⣿⣿⣿⣿⣿⣿⣿⣦⣄⠀⢀⣿⣿⣿⣿⡿⠿⠚⠀⢀⣾⣧⡀⠀⠀⠀⠀⠑⠄⠀
⠀⠀⠀⣤⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣤⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣜⠇⣇⠹⠋⠀⠀⠀⠀⢨⣿⠿⢦⡀⠀⠀⠀⠀⡀⠀
⠀⠀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⣿⣏⠀⠀⠀⠀⣀⣿⠇⠀⠀⢷⣆⣸⣀⣰⣿⣇
⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣸⣿⣿⣷⣶⡾⢟⠁⠀⠀⢀⠀⣹⣟⡉⠀⠈⠻
⢀⣿⡏⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠉⠉⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⣾⣤⣤⣴⣾⣿⣿⣾⡆⠀⠀⠀
⠀⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⣿⣿⣿⠿⣿⣿⠁⠀⠀⠀⠀⢸⣿⣿⣿⠇⠀⠀⠙⣿⠛⢙⣛⠿⢿⣿⣿⣿⡯⠿⠿⠏⠏⢃⠀⠀⠀
⠀⢠⡏⠁⠀⠀⠀⠙⣻⣿⣿⣿⡿⠿⠻⠿⡛⡀⣾⣿⣦⣤⡄⣀⣀⣰⣿⠟⣿⣂⠀⠀⠀⣷⠀⣺⣿⣧⢈⣥⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢀⡆⢄⡄⠀⢀⣠⣴⣿⣿⣿⠿⠀⠀⠀⠀⠈⢣⣿⣿⣿⣿⣷⣿⣿⣿⣿⠀⢘⣿⣷⣶⣾⣿⡆⠸⠿⢿⣿⠟⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢿⢿⣿⣷⣴⣿⣿⣿⠟⠁⢸⣧⡄⠀⠀⣀⣴⡷⠻⢿⣿⣿⣿⣿⣿⣿⣧⣠⣤⣿⣿⣿⣿⣿⣿⠀⠀⠈⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢨⢻⣿⣿⣿⣿⣿⡃⢀⣀⢀⣯⣿⣿⣶⣿⣿⡇⠀⠀⠉⠉⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠟⠻⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢸⣿⣿⣿⣿⣿⣾⣯⣴⣿⣿⢿⣿⣿⡿⠃⠀⠀⠀⠘⢾⣿⡟⡿⠛⠛⠛⢿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠐⠿⠿⣿⢿⡿⢿⣿⣿⣿⣿⣽⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠁⠀⠀⠀⠀⠀⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠛⠛⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

    DNS SERVER IS LISTENING...
`

func parseQuesSection(buf []byte, start int) (string, int) {
	domain := ""
	curr := start

	for {
		length := int(buf[curr]) // length byte
		if length == 0 {         //null label
			curr++
			break
		}

		curr++ //moves to label

		label := string(buf[curr : curr+length])

		if domain == "" {
			domain = label
		} else {
			domain += "." + label
		}

		curr += length
	}

	return domain, curr
}

func main() {

	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")} //dia chi nghe
	//mo cong ket noi bang listenUDP
	connect, errr := net.ListenUDP("udp", &addr) //tra ve UDPConn va error

	if errr != nil {
		fmt.Println("Error:", errr)
		return
	}

	defer connect.Close() //phong truong hop bi loi thi van co the out

	fmt.Println(banner)

	for {
		//chuan bi nhan du lieu
		buf := make([]byte, 512)
		//ReadFromUDP tra ve int(number of byte), Addr va error
		n, address, err := connect.ReadFromUDP(buf)

		//if it not say anything try "nslookup google.com 127.0.0.1" in another powershell
		fmt.Printf("Received %d bytes from %s:\n", n, address)

		//have fun with header
		var id, qr, opcode, rd, qdCount uint16
		var msgType string = "Unknown"
		if n >= 12 {
			fmt.Println("HEADER SECTION")

			id = binary.BigEndian.Uint16(buf[0:2])

			flags := binary.BigEndian.Uint16(buf[2:4])
			//extracting bits
			qr = flags >> 15
			if qr == 1 {
				msgType = "Response"
			} else {
				msgType = "Query"
			}
			opcode = (flags >> 11) & 0xF
			rd = (flags >> 8) & 1

			qdCount = binary.BigEndian.Uint16(buf[4:6])

			/*anCount := binary.BigEndian.Uint16(buf[6:8])

			  nsCount := binary.BigEndian.Uint16(buf[8:10])

			  arCount := binary.BigEndian.Uint16(buf[10:12])

			  so the question section should be from 12 to 17??
			*/
			fmt.Printf("ID: %v\t", id)
			fmt.Printf("QR: %v\t", msgType)
			fmt.Printf("Opcode: %d\t", opcode)
			fmt.Printf("RD: %d\t", rd)
			fmt.Printf("QDCOUNT: %d\n", qdCount)
		}

		fmt.Println("QUESTION SECTION")

		domain, curr := parseQuesSection(buf, 12)
		//qtype will start after null label (2 bytes), qclass 2 bytes
		qType := binary.BigEndian.Uint16(buf[curr : curr+2])
		qClass := binary.BigEndian.Uint16(buf[curr+2 : curr+4])

		fmt.Printf("QNAME: %s\t", domain)
		fmt.Printf("QTYPE: %d\t", qType)
		fmt.Printf("QCLASS: %d\n", qClass)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
	}

}
