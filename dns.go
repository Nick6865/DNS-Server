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

var blacklistMap map[string]bool

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16 // Number of Questions
	ANCount uint16 // Number of Answers
	NSCount uint16 // Number of Authorities
	ARCount uint16 // Number of Additionals
}

type DNSQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

func ParseHeader(buf []byte) DNSHeader {
	return DNSHeader{
		ID:      binary.BigEndian.Uint16(buf[0:2]),
		Flags:   binary.BigEndian.Uint16(buf[2:4]),
		QDCount: binary.BigEndian.Uint16(buf[4:6]),
		ANCount: binary.BigEndian.Uint16(buf[6:8]),
		NSCount: binary.BigEndian.Uint16(buf[8:10]),
		ARCount: binary.BigEndian.Uint16(buf[10:12]),
	}
}

func (h *DNSHeader) Pack(buf []byte) {
	binary.BigEndian.PutUint16(buf[0:2], h.ID)
	binary.BigEndian.PutUint16(buf[2:4], h.Flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)
}

func (h *DNSHeader) SetNXDOMAIN() {
	h.Flags |= (1 << 15) // QR = 1 (Response)
	h.Flags |= (1 << 10) // AA = 1 (Authoritative)
	h.Flags |= (1 << 7)  // RA = 1 (Recursion Available)
	h.Flags &= 0xFFF0    // delete old rdcode
	h.Flags |= 3         // RCODE = 3 (NXDOMAIN)
}

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

func handlePacket(connect *net.UDPConn, upstream net.Conn, address net.Addr, n int, buf []byte, blacklist map[string]bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf(" [!] Packet error: %v\n", r)
		}
	}()

	if n < 12 {
		return
	}

	header := ParseHeader(buf)
	now := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("\n\033[36m[%s] Received %d bytes from %s\033[0m\n", now, n, address.String())
	fmt.Printf("\033[33mHEADER SECTION:\033[0m\n")
	fmt.Printf(" [Log] Query ID: %d | Questions: %d\n", header.ID, header.QDCount)

	curr := 12

	fmt.Printf("\033[33mQUESTION SECTION:\033[0m\n")

	for i := 0; i < int(header.QDCount); i++ {
		domain, nextPos, err := parseName(buf, curr, n)
		if err != nil {
			return
		}

		//checking blacklist
		if blacklist[domain] {
			fmt.Printf(" \033[31m[!] Domain blocked: %s\033[0m\n", domain)

			//use struct to make it faster
			header.SetNXDOMAIN()

			reply := make([]byte, n)
			copy(reply, buf[:n])
			header.Pack(reply) //write fixed header to packet

			udpAddr, _ := address.(*net.UDPAddr)
			connect.WriteToUDP(reply, udpAddr)
			return
		}

		qType := binary.BigEndian.Uint16(buf[nextPos : nextPos+2])
		qClass := binary.BigEndian.Uint16(buf[nextPos+2 : nextPos+4])

		fmt.Printf(" [Log] Requesting: \033[36m%s\033[0m\t| Type: %d\t| Class: %d\n", domain, qType, qClass)

		curr = nextPos + 4
	}

	// connect to upstream and add timeout
	upstream, err := net.DialTimeout("udp", "8.8.8.8:53", 1*time.Second)
	if err != nil {
		fmt.Println(" [!] Upstream error:", err)
		return
	}
	defer upstream.Close()

	//send request
	_, err = upstream.Write(buf[:n])
	if err != nil {
		fmt.Println(" [!] Failed to forward query:", err)
		return
	}

	//set deadline if >1sec skip it
	upstream.SetReadDeadline(time.Now().Add(1 * time.Second))

	//receive answer
	reply := make([]byte, 1024)
	nReply, err := upstream.Read(reply)
	if err != nil {
		fmt.Println(" [!] No response from Google (Timeout/Error):", err)
		return
	}

	resID := binary.BigEndian.Uint16(reply[0:2])

	fmt.Printf("\033[32m[RESOLVED]\033[0m | ID: %05d | Bytes: %d\n", resID, nReply)

	udpAddr, ok := address.(*net.UDPAddr)
	if !ok {
		fmt.Println(" [!] Invalid client address type")
		return
	}

	_, err = connect.WriteToUDP(reply[:nReply], udpAddr)
	if err != nil {
		fmt.Println(" [!] Failed back to client:", err)
	} else {
		fmt.Printf(" \033[32m[OK] Answered %s\033[0m\n", udpAddr.String())
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
		go handlePacket(connect, upstream, ClientAddress, n, packetData, blacklistMap)
	}
}

//if it not say anything try "nslookup google.com 127.0.0.1" in another powershell
