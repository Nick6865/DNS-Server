package main

import (
	"encoding/binary"
	"fmt"
)

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
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

	h.ANCount = 0
	h.NSCount = 0
	h.ARCount = 0
}

func parseName(buf []byte, curr int, maxLen int) (string, int, error) {
	ptr := curr
	domain := ""
	jumped := false
	nextPos := -1
	jumpCount := 0
	const maxJump = 5

	for {
		if jumpCount > maxJump {
			return "", 0, fmt.Errorf("too many jumps (potential loop)")
		}
		if ptr >= maxLen {
			return "", 0, fmt.Errorf("out of bounds")
		}

		lenByte := buf[ptr]

		if (lenByte & 0xC0) == 0xC0 {
			if ptr+1 >= maxLen {
				return "", 0, fmt.Errorf("truncated pointer")
			}
			offset := int(binary.BigEndian.Uint16(buf[ptr:ptr+2]) & 0x3FFF)
			if !jumped {
				nextPos = ptr + 2
			}
			ptr = offset
			jumped = true
			jumpCount++
			continue
		}

		ptr++
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
