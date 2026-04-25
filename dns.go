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

func parseName(buf []byte, curr int, maxLen int) (string, int, error) {
	ptr := curr //will jump if meet compression
	domain := ""
	jumped := false
	nextPos := -1 // pos to read qtype after qname
	jumpCount := 0
	const maxJump = 5 //avoid loops

	for {
		if jumpCount > maxJump {
			return "", 0, fmt.Errorf("too many jumps (potential loop)")
		}
		if ptr >= maxLen {
			return "", 0, fmt.Errorf("out of bounds")
		}

		lenByte := buf[ptr]

		// checking first 2 bits (11 is a pointer)
		if (lenByte & 0xC0) == 0xC0 {
			if ptr+1 >= maxLen {
				return "", 0, fmt.Errorf("truncated pointer")
			}

			// 14 bits offset
			offset := int(binary.BigEndian.Uint16(buf[ptr:ptr+2]) & 0x3FFF) //check where is the pointer point to

			if !jumped {
				nextPos = ptr + 2
			}

			ptr = offset //go back to that address
			jumped = true
			jumpCount++
			continue
		}

		// read normal label
		ptr++ //skip the length bit
		if lenByte == 0 {
			break
		}

		if ptr+int(lenByte) > maxLen {
			return "", 0, fmt.Errorf("label length out of range")
		}

		label := string(buf[ptr : ptr+int(lenByte)])
		if domain == "" {
			domain = label
		} else {
			domain += "." + label
		}
		ptr += int(lenByte)
	}

	if !jumped {
		nextPos = ptr
	}

	return domain, nextPos, nil
}

func main() {

	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")} //dia chi nghe
	//mo cong ket noi bang listenUDP
	connect, errr := net.ListenUDP("udp", &addr) //tra ve UDPConn va error

	//lets connect to google and become a fowarder
	network, dialErr := net.Dial("udp", "8.8.8.8:53") //its google public dns server port 53

	if dialErr != nil {
		fmt.Println("Error: ", dialErr)
		return
	}
	defer network.Close()

	if errr != nil {
		fmt.Println("Error:", errr)
		return
	}

	defer connect.Close() //phong truong hop bi loi thi van co the out

	fmt.Println(banner2)

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

			fmt.Println("QUESTION SECTION")
			fmt.Printf("--- Parsing %d Question(s) ---\n", qdCount)

			curr := 12

			for i := 0; i < int(qdCount); i++ {
				//parsing the name
				domain, nextPos, err := parseName(buf, curr, n)
				if err != nil {
					fmt.Printf("Error parsing question %d: %v\n", i, err)
					break
				}

				// check next 4 byte(qtype and qclass)
				if nextPos+4 > n {
					fmt.Println("Error: Not enough data for QTYPE/QCLASS")
					break
				}

				qType := binary.BigEndian.Uint16(buf[nextPos : nextPos+2])
				qClass := binary.BigEndian.Uint16(buf[nextPos+2 : nextPos+4])

				fmt.Printf("#%d | Name: %s | Type: %d | Class: %d\n", i+1, domain, qType, qClass)

				// update curr
				curr = nextPos + 4
			}
		}
		//send question section and receive a reply
		network.Write(buf[:n]) //sending
		reply := make([]byte, 512)
		nreply, rError := network.Read(reply)

		if rError != nil {
			fmt.Println("Error: ", rError)
			continue
		}

		_, wError := connect.WriteToUDP(reply[:nreply], address)

		if wError != nil {
			fmt.Println("Error: ", wError)
			continue
		}

		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
	}

}

//if it not say anything try "nslookup google.com 127.0.0.1" in another powershell
