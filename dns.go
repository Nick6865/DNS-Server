package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
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

type ServerStats struct {
	mu            sync.Mutex
	TotalRequests int
	Blocked       int
	CacheHits     int
}

var stats ServerStats
var logger *log.Logger

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func runCLIDashboard() {
	for {
		clearScreen()

		stats.mu.Lock()
		total := stats.TotalRequests
		blocked := stats.Blocked
		hits := stats.CacheHits
		stats.mu.Unlock()

		cacheRate := 0.0
		if total > 0 {
			cacheRate = float64(hits) / float64(total) * 100
		}
		//drawing/designing
		fmt.Println(banner)

		fmt.Println("\033[36m==================================================\033[0m")
		fmt.Println("\033[1;32m           DNS FORWARDER CLI DASHBOARD          \033[0m")
		fmt.Println("\033[36m==================================================\033[0m")
		fmt.Printf(" \033[1;37m Total hits:        \033[0m %8d requests\n", total)
		fmt.Printf(" \033[1;31m Blocked (Ads/Bad): \033[0m %8d requests\n", blocked)
		fmt.Printf(" \033[1;33m Cache Hits (Yes):  \033[0m %8d requests\n", hits)
		fmt.Printf(" \033[1;35m Cache ratio:       \033[0m %8.2f %%\n", cacheRate)
		fmt.Println("\033[36m==================================================\033[0m")
		fmt.Println(" Access details are being logged to the file: dns.log")
		fmt.Println(" Ctrl + C to turn off server...")

		//rest
		time.Sleep(20 * time.Second)
	}
}

func handlePacket(connect *net.UDPConn, upstream net.Conn, address net.Addr, n int, buf []byte, blacklist map[string]bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf(" [!] Packet error: %v\n", r)
		}
	}()

	stats.mu.Lock()
	stats.TotalRequests++
	stats.mu.Unlock()

	if n < 12 {
		return
	}

	header := ParseHeader(buf)
	//now := time.Now().Format("2006-01-02 15:04:05")

	//fmt.Printf("\n\033[36m[%s] Received %d bytes from %s\033[0m\n", now, n, address.String())
	//fmt.Printf(" [Log] Query ID: %d | Questions: %d\n", header.ID, header.QDCount)

	curr := 12
	var cacheKey string
	var domainN string

	for i := 0; i < int(header.QDCount); i++ {
		domain, nextPos, err := parseName(buf, curr, n)
		domainN = domain
		if err != nil {
			return
		}

		if blacklist[domain] {
			stats.mu.Lock()
			stats.Blocked++
			stats.mu.Unlock()
			logger.Printf("[BLOCKED] %s", domain)

			//fmt.Printf(" \033[31m[!] Domain blocked: %s\033[0m\n", domain)
			header.SetNXDOMAIN()
			reply := make([]byte, n)
			copy(reply, buf[:n])
			header.Pack(reply)
			udpAddr, _ := address.(*net.UDPAddr)
			connect.WriteToUDP(reply, udpAddr)
			return
		}

		qType := binary.BigEndian.Uint16(buf[nextPos : nextPos+2])
		//qClass := binary.BigEndian.Uint16(buf[nextPos+2 : nextPos+4])

		//fmt.Printf(" [Log] Requesting: \033[36m%s\033[0m\t| Type: %d\t| Class: %d\n", domain, qType, qClass)

		cacheKey = fmt.Sprintf("%s|%d", domain, qType)

		memoryCache.mu.RLock()
		entry, found := memoryCache.store[cacheKey]
		memoryCache.mu.RUnlock()

		if found && time.Now().Before(entry.ExpireAt) {
			stats.mu.Lock()
			stats.CacheHits++
			stats.mu.Unlock()
			logger.Printf("[CACHE HIT] %s", domain)

			//fmt.Printf(" \033[35m[CACHE HIT]\033[0m Return %s immediately!\n", domain)
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
		//fmt.Println(" [!] Failed to forward query:", err)
		return
	}

	upstream.SetReadDeadline(time.Now().Add(1 * time.Second))

	reply := make([]byte, 4096)
	nReply, err := upstream.Read(reply)
	if err != nil {
		//fmt.Println(" [!] No response from Google (Timeout/Error):", err)
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

	//resID := binary.BigEndian.Uint16(reply[0:2])

	logger.Printf("[RESOLVED] %s (from Google)", domainN)

	//fmt.Printf("\033[32m[RESOLVED]\033[0m | ID: %05d | Bytes: %d\n", resID, nReply)

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

	fileLog, err := os.OpenFile("dns.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Cannot create log:", err)
	}
	defer fileLog.Close()
	logger = log.New(fileLog, "", log.Ldate|log.Ltime)
	logger.Println("=== ACTIVATING SERVER ===")

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

	go runCLIDashboard()

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
