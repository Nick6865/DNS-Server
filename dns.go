package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const banner = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣶⣶⣶⣶⣶⣶⣦⣄⡀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⢀⣀⣀⣀⡀⠀⠀⠀⠀⠀⢀⣈⣿⣶⣿⣿⣿⡿⢿⣿⡿⠛⠙⠻⣿⣿⡇⠀⠀
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
	fmt.Printf(" [Log] Query ID: %d | Questions: %d\n", header.ID, header.QDCount)

	curr := 12
	var cacheKey string

	for i := 0; i < int(header.QDCount); i++ {
		domain, nextPos, err := parseName(buf, curr, n)
		if err != nil {
			return
		}

		if blacklist[domain] {
			fmt.Printf(" \033[31m[!] Domain blocked: %s\033[0m\n", domain)
			header.SetNXDOMAIN()
			reply := make([]byte, n)
			copy(reply, buf[:n])
			header.Pack(reply)
			udpAddr, _ := address.(*net.UDPAddr)
			connect.WriteToUDP(reply, udpAddr)
			return
		}

		qType := binary.BigEndian.Uint16(buf[nextPos : nextPos+2])
		qClass := binary.BigEndian.Uint16(buf[nextPos+2 : nextPos+4])

		fmt.Printf(" [Log] Requesting: \033[36m%s\033[0m\t| Type: %d\t| Class: %d\n", domain, qType, qClass)

		cacheKey = fmt.Sprintf("%s|%d", domain, qType)

		memoryCache.mu.RLock()
		entry, found := memoryCache.store[cacheKey]
		memoryCache.mu.RUnlock()

		if found && time.Now().Before(entry.ExpireAt) {
			fmt.Printf(" \033[35m[CACHE HIT]\033[0m Return %s immediately!\n", domain)
			cachedReply := make([]byte, len(entry.Response))
			copy(cachedReply, entry.Response)
			binary.BigEndian.PutUint16(cachedReply[0:2], header.ID)
			udpAddr, _ := address.(*net.UDPAddr)
			connect.WriteToUDP(cachedReply, udpAddr)
			return
		}

		curr = nextPos + 4
	}

	_, err := upstream.Write(buf[:n])
	if err != nil {
		fmt.Println(" [!] Failed to forward query:", err)
		return
	}

	upstream.SetReadDeadline(time.Now().Add(1 * time.Second))

	reply := make([]byte, 4096)
	nReply, err := upstream.Read(reply)
	if err != nil {
		fmt.Println(" [!] No response from Google (Timeout/Error):", err)
		return
	}

	if cacheKey != "" {
		memoryCache.mu.Lock()
		memoryCache.store[cacheKey] = CacheEntry{
			Response: reply[:nReply],
			ExpireAt: time.Now().Add(5 * time.Minute),
		}
		memoryCache.mu.Unlock()
	}

	resID := binary.BigEndian.Uint16(reply[0:2])
	fmt.Printf("\033[32m[RESOLVED]\033[0m | ID: %05d | Bytes: %d\n", resID, nReply)

	udpAddr, ok := address.(*net.UDPAddr)
	if !ok {
		return
	}

	connect.WriteToUDP(reply[:nReply], udpAddr)
}

func main() {
	var err error
	blacklistMap, err = loadBlackList("blacklist.txt")

	if err != nil {
		fmt.Println("Warning: cannot find blacklist.txt, creating a new blacklistMap")
		blacklistMap = make(map[string]bool)
	} else {
		fmt.Printf("Add %d domain to blacklist.\n", len(blacklistMap))
	}

	addr := net.UDPAddr{Port: 53, IP: net.ParseIP("0.0.0.0")}
	connect, errr := net.ListenUDP("udp", &addr)
	if errr != nil {
		fmt.Println("Error:", errr)
		return
	}
	defer connect.Close()

	upstream, err := net.DialTimeout("udp", "8.8.8.8:53", 2*time.Second)
	if err != nil {
		fmt.Println("Failed to connect to upstream:", err)
		return
	}
	defer upstream.Close()

	fmt.Println(banner)
	fmt.Println("DNS FORWARDER IS RUNNING ON PORT 53...")

	for {
		buf := make([]byte, 4096)
		n, ClientAddress, err := connect.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		packetData := make([]byte, n)
		copy(packetData, buf[:n])

		go handlePacket(connect, upstream, ClientAddress, n, packetData, blacklistMap)
	}
}
