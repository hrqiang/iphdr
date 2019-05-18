package main

import (
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}

func pkt() []byte {
	h := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + 10, // 20 bytes for IP, 10 for ICMP
		TTL:      64,
		Protocol: 1, // ICMP
		Dst:      net.IPv4(172, 217, 164, 100),
		// ID, Src and Checksum will be set for us by the kernel
	}

	icmp := []byte{
		8, // type: echo request
		0, // code: not used by echo request
		0, // checksum (16 bit), we fill in below
		0,
		0, // identifier (16 bit). zero allowed.
		0,
		0, // sequence number (16 bit). zero allowed.
		0,
		0xC0, // Optional data. ping puts time packet sent here
		0xDE,
	}
	cs := csum(icmp)
	icmp[2] = byte(cs)
	icmp[3] = byte(cs >> 8)

	out, err := h.Marshal()
	fmt.Printf("IP header: % X\n", out)
	if err != nil {
		log.Fatal(err)
	}
	return append(out, icmp...)
}

func main() {
	go receiveAndPrint()

	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{172, 217, 164, 100},
	}
	p := pkt()
	err := syscall.Sendto(fd, p, 0, &addr)
	if err != nil {
		log.Fatal("Sendto:", err)
	}

	time.Sleep(time.Second * 5)
}

func receiveAndPrint() {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Ping reply from Google % X\n", buf[:numRead])
	}
}
