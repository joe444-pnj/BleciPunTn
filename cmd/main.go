package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"IP-Info-Tool/pkg/ipinfo"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func main() {
	displayAsciiArt()
	color.Cyan("Welcome to the Enhanced IP Info Tool!")
	color.Cyan("Type 'help' for a list of commands.")
	color.Cyan("Enter IP addresses (comma-separated for multiple IPs) or type a command:")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		switch strings.ToLower(input) {
		case "exit":
			color.Yellow("Exiting the tool. Goodbye!")
			return
		case "clear":
			clearScreen()
		case "help":
			displayHelp()
		default:
			ipList := strings.Split(input, ",")
			runTool(ipList)
		}
	}
}

func displayAsciiArt() {
	asciiArt := `
	  ____  _           _ _____           _____       
	 |  _ \| |         (_)  __ \         |_   _|      
	 | |_) | | ___  ___ _| |__) |   _ _ __ | |  _ __  
	 |  _ <| |/ _ \/ __| |  ___/ | | | '_ \| | | '_ \ 
	 | |_) | |  __/ (__| | |   | |_| | | | | |_| | | |
	 |____/|_|\___|\___|_|_|    \__,_|_| |_|____|_| |_|
	`
	color.Green(asciiArt)
}

func displayHelp() {
	color.Cyan("Commands available:")
	color.Cyan("  clear  - Clears the screen")
	color.Cyan("  exit   - Exits the tool")
	color.Cyan("  help   - Displays this help message")
	color.Cyan("Usage:")
	color.Cyan("  Enter one or multiple IP addresses separated by commas to scan them.")
	color.Cyan("  Example: 192.168.1.1,8.8.8.8")
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	color.Cyan("Screen cleared. Enter new IP addresses or type a command:")
}

func runTool(ipList []string) {
	var wg sync.WaitGroup
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()

	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			ip = strings.TrimSpace(ip)
			if !validateIP(ip) {
				color.Red("Invalid IP address: %s", ip)
				return
			}
			ipInfo, err := ipinfo.FetchIPInfo(ip)
			if err != nil {
				color.Red("Error fetching IP info for %s: %s", ip, err)
				return
			}
			s.Stop()
			displayIPInfo(ipInfo)
			color.Blue("Scanning ports for %s...", ipInfo.IP)
			scanPorts(ipInfo.IP, 1, 1024) // Default port range
		}(ip)
	}

	wg.Wait()
	s.Stop()
	color.Green("Scanning completed.")
	logScanResults(ipList)
}

func validateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func displayIPInfo(ipInfo *ipinfo.IPInfo) {
	color.Cyan("IP: %s", ipInfo.IP)
	color.Green("Hostname: %s", ipInfo.Hostname)
	color.Yellow("Location: %s, %s, %s", ipInfo.City, ipInfo.Region, ipInfo.Country)
	color.Red("Timezone: %s", ipInfo.Timezone)
	color.Cyan("Organization: %s", ipInfo.Org)
	color.Magenta("Reverse DNS: %s", reverseDNSLookup(ipInfo.IP))
}

func reverseDNSLookup(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "N/A"
	}
	return strings.Join(names, ", ")
}

func scanPorts(ip string, startPort, endPort int) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 20) // Limit concurrency to 20 goroutines

	// Iterate through the range of ports
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			address := fmt.Sprintf("%s:%d", ip, port)
			// Attempt to connect to the port with a timeout of 2 seconds
			conn, err := net.DialTimeout("tcp", address, 2*time.Second)
			if err == nil {
				color.Green("Port %d open on %s - Service: %s", port, ip, identifyService(port))
				bannerGrab(conn)
				conn.Close()
			} else {
				color.Red("Error connecting to %s:%d - %s", ip, port, err)
			}
		}(port)
	}

	wg.Wait()
	color.Green("Port scan completed for %s", ip)
}

func identifyService(port int) string {
	services := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		443:  "HTTPS",
		445:  "Microsoft-DS",
		3306: "MySQL",
		3389: "RDP",
	}

	if service, found := services[port]; found {
		return service
	}
	return "Unknown"
}

func bannerGrab(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == nil {
		color.Magenta("Banner: %s", strings.TrimSpace(string(buf[:n])))
	}
}

func logScanResults(ipList []string) {
	logsDir := "logs"
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		os.Mkdir(logsDir, os.ModePerm)
	}

	filePath := filepath.Join(logsDir, "scan_results.log")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("Error opening log file: %s", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	writer.WriteString(fmt.Sprintf("Scan results at %s:\n", timestamp))

	for _, ip := range ipList {
		writer.WriteString(fmt.Sprintf("- %s\n", ip))
	}
	writer.WriteString("\n")
}
