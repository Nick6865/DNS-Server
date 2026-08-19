# Local DNS Ad-Blocker

A lightweight local DNS server written in Go that blocks ads and 
trackers, similar in concept to Pi-hole. It resolves known 
ad-serving/tracking domains to a null address instead of forwarding 
them to the real internet, and forwards all other queries to a 
normal upstream DNS resolver.

## Overview

This project was built as a learning exercise to understand the DNS 
protocol at the wire format level — parsing raw DNS headers and 
question sections instead of relying on high-level DNS libraries.

## Features

- **DNS query parsing**: decodes DNS header and question section 
  from raw UDP packets
- **Blocklist-based filtering**: checks each requested domain against 
  a local blocklist (`blacklist.txt`)
- **Ad/tracker blocking**: returns a null response (0.0.0.0) for 
  blocklisted domains instead of forwarding them
- **Query logging**: logs incoming queries for debugging and 
  inspection during development
- **Caching**: basic response caching to reduce repeated lookups

## How it works

1. Listens for DNS queries on UDP port 53 on the local network
2. Parses the DNS header and question section
3. Checks the requested domain against the blocklist
4. If blocked → returns 0.0.0.0 (ad/tracker blocked)
5. If not blocked → forwards the query to a real upstream resolver 
   and returns the actual response

## Installation

```bash
git clone https://github.com/Nick6865/DNS-Server
cd DNS-Server
go build dns.go
```

## Usage

```bash
sudo ./dns
```

## Notes on logging

Query logs (`dns.log`) can grow quickly during testing. This is a 
learning/demo project, not a production log-management setup — for 
now, logs are cleared on shutdown (Ctrl+C) simply to avoid 
accumulating large files during repeated test runs on a personal 
machine. This is **not** intended as a way to conceal activity; 
in a real deployment, log rotation (e.g. capping file size and 
archiving old entries) would replace this, since permanently 
deleting logs is not appropriate for any tool meant to run 
persistently or be relied on for troubleshooting.

## Scope

This project only intercepts and responds to DNS queries on the 
local network for ad-blocking purposes — it does not spoof or 
redirect traffic for any domain outside the blocklist, and it is 
only intended to run on networks/devices the user owns or has 
explicit permission to configure.

## Disclaimer

Built for personal learning about the DNS protocol and network 
programming in Go. Not intended for production use as-is.
