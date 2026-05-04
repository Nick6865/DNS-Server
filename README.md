# DNS Sniffer & Redirection Framework
A high-performance DNS interception and redirection tool written in Go, designed for network security auditing and traffic analysis. This framework enables real-time monitoring of DNS queries and conditional redirection based on predefined rules.

## Overview
This project provides a robust solution for capturing DNS traffic at the packet level. By leveraging Go's concurrency model, it ensures low-latency processing even under significant network load.

## Key Features
- **Packet-Level Interception:** Utilizes raw sockets to monitor and parse DNS queries.
- **Dynamic Redirection:** Implements logic to redirect specific domains to designated IP addresses.
- **Comprehensive Logging:** Records all intercepted queries and responses for forensic analysis.
- **Integrated Kill Switch:** A built-in security mechanism to immediately terminate processes and purge sensitive logs upon command.
- **Stateless Operation:** Designed for minimal memory footprint and high stability.

## Technical Architecture
The framework operates by binding to port 53 and analyzing incoming UDP packets. It decodes the DNS wire format to extract queries and applies a redirection filter:
1. **Interception:** Captures incoming UDP packets on the local interface.
2. **Parsing:** Extracts the Domain Name System (DNS) header and question section.
3. **Filtering:** Checks the requested domain against a local redirection table.
4. **Response:** Forges a DNS response (A Record) pointing to the target IP if a match is found.

## Installation
Ensure you have the Go runtime installed on your system.
```bash
# Clone the repository
git clone [https://github.com/Nick6865/DNS-Server](https://github.com/Nick6865/DNS-Server)
# Navigate to the project directory
cd DNS-Server
# Build the binary
go build dns.go
```

## Usage
```bash
sudo ./dns
```

## Safety Mechanism (Kill Switch)
To trigger the automated log purge and safe shutdown, use the designated interrupt signal (Ctrl+C). The killSwitch() function will execute immediately to ensure no data remains in the active log files.

## Legal Disclaimer
This software is provided for educational purposes and authorized security testing only. Unauthorized use of this tool against networks without explicit permission is illegal and unethical. The developer assumes no liability for misuse or damage caused by this software.
