package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var blacklistMap map[string]bool

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

func loadBlackList(filename string) (map[string]bool, error) {
	list := make(map[string]bool)
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" && !strings.HasPrefix(domain, "#") {
			list[domain] = true
		}
	}
	return list, scanner.Err()
}

func sendNXDomain(connect *net.UDPConn, address net.Addr, buf []byte) {
	reply := make([]byte, len(buf))
	copy(reply, buf)
	//change the packet(flags in header)
	/*reply[2] = 0x81 //10000001
	reply[3] = 0x83 //10000011 the last 4 bit is NXDOMAIN error*/

	flags := binary.BigEndian.Uint16(buf[2:4])
	var responseFlags uint16

	responseFlags |= (1 << 15)          //Bit QR (15): 0 -> 1 (query to response)
	responseFlags |= (1 << 10)          //Bit AA (10): 1 (our server holding this question)
	responseFlags |= (flags & (1 << 8)) //RD keeps this RD from client
	responseFlags |= (1 << 7)           //RA (Server supports recursion)
	responseFlags |= 3                  //RCODE: NXDOMAIN in last 4 bits

	binary.BigEndian.PutUint16(reply[2:4], responseFlags)

	udpAddr, ok := address.(*net.UDPAddr)
	if !ok {
		fmt.Println("Invalid client address type")
		return
	}
	_, err := connect.WriteToUDP(reply, udpAddr)
	if err != nil {
		fmt.Println("Failed back to client:", err)
	} else {
		fmt.Printf(" [OK] Answered %s\n", udpAddr.String())
	}

	fmt.Printf(" [!] Blocked & Sent NXDOMAIN to %s\n", udpAddr.String())
}

func handlePacket(connect *net.UDPConn, upstream net.Conn, address net.Addr, n int, buf []byte, blacklist map[string]bool) {
	defer recover()
	// using parse to get query info to print
	if n >= 12 {
		id := binary.BigEndian.Uint16(buf[0:2])
		qdCount := binary.BigEndian.Uint16(buf[4:6])

		fmt.Print("HEADER SECTION: \n")
		fmt.Printf(" [Log] Query ID: %d | Questions: %d\n", id, qdCount)

		curr := 12
		for i := 0; i < int(qdCount); i++ {
			domain, nextPos, err := parseName(buf, curr, n)

			//if the domain is an ad domain or you don't want this domain to appear create a blacklist to block it (send a fake error packet to google by using bitwise manipulation)
			if blacklist[domain] {
				fmt.Printf(" [!] Domain blocked: %s\n", domain)
				sendNXDomain(connect, address, buf)
				return
			}

			if err == nil {
				// check next 4 byte(qtype and qclass)
				if nextPos+4 > n {
					fmt.Println("Error: Not enough data for QTYPE/QCLASS")
					break
				}
				fmt.Print("QUESTION SECTION: \n")
				qType := binary.BigEndian.Uint16(buf[nextPos : nextPos+2])
				qClass := binary.BigEndian.Uint16(buf[nextPos+2 : nextPos+4])
				fmt.Printf(" [Log] Requesting: %s \t| Type: %d | Class: %d\n", domain, qType, qClass)
				curr = nextPos + 4
			}
		}
	}

	_, err := upstream.Write(buf[:n])
	if err != nil {
		fmt.Println("Failed to forward query:", err)
		return
	}

	// receive a reply
	reply := make([]byte, 2048)
	//set deadline
	nReply, err := upstream.Read(reply)
	if err != nil {
		fmt.Println("Failed to read reply from upstream:", err)
		return
	}

	//send back to client
	udpAddr, ok := address.(*net.UDPAddr)
	if !ok {
		fmt.Println("Invalid client address type")
		return
	}

	_, err = connect.WriteToUDP(reply[:nReply], udpAddr)
	if err != nil {
		fmt.Println("Failed back to client:", err)
	} else {
		fmt.Printf(" [OK] Answered %s\n", udpAddr.String())
	}
}

func main() {
	//creating a blacklist
	var err error
	blacklistMap, err = loadBlackList("blacklist.txt")

	if err != nil {
		fmt.Println("Warning: cannot fine blacklist.txt, creating a new blacklistMap")
		blacklistMap = make(map[string]bool)
	} else {
		fmt.Printf("Add %d domain to blacklist.\n", len(blacklistMap))
	}

	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")} //dia chi nghe
	//mo cong ket noi bang listenUDP
	connect, errr := net.ListenUDP("udp", &addr) //tra ve UDPConn va error

	//forward to google dns
	upstream, err := net.DialTimeout("udp", "8.8.8.8:53", 2*time.Second)
	if err != nil {
		fmt.Println("Failed to connect to upstream:", err)
		return
	}
	defer upstream.Close()

	if errr != nil {
		fmt.Println("Error:", errr)
		return
	}
	defer connect.Close() //phong truong hop bi loi thi van co the out

	fmt.Println(banner)

	fmt.Println("DNS FORWARDER IS RUNNING ON PORT 53...")

	for {
		//add time
		now := time.Now().Format("2006-01-02 15:04:05")

		//chuan bi nhan du lieu
		buf := make([]byte, 512)
		//ReadFromUDP tra ve int(number of byte), Addr va error
		n, ClientAddress, err := connect.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Read Error:", err)
			continue
		}

		packetData := make([]byte, n)
		copy(packetData, buf[:n])

		//if it not say anything try "nslookup google.com 127.0.0.1" in another powershell
		fmt.Printf("\n[%s] Received %d bytes from %s:\n", now, n, ClientAddress)

		go handlePacket(connect, upstream, ClientAddress, n, packetData, blacklistMap)
	}
}

//if it not say anything try "nslookup google.com 127.0.0.1" in another powershell
