package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

//const r_str = "GET %s HTTP/1.1\r\nHost: %s\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\r\nAccept-Encoding: gzip, deflate, br\r\nAccept-Language: ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7\r\n\r\n"

const r_str = "GET %s HTTP/1.1\r\n\r\n"

func request_direct(url, addr string, port int, timeout int) ([]byte, bool) {
	ap := fmt.Sprintf("%s:%d", addr, port)
	//request_bytes := []byte(fmt.Sprintf(r_str, url, ap))
	request_bytes := []byte(fmt.Sprintf(r_str, url))
	conn, err := net.Dial("tcp", ap)
	if err != nil {
		// connection error
		fmt.Println("connection error", err)
		return nil, false
	}

	defer conn.Close()
	conn.Write(request_bytes)
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 4096)    // using small tmo buffer for demonstrating
	duration := time.Duration(timeout) * time.Second

	for {
		if timeout > 0 {
			err := conn.SetReadDeadline(time.Now().Add(duration))
			if err != nil {
				log.Println("SetReadDeadline failed:", err)
				// do something else, for example create new conn
				return nil, false
			}
		}
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		//fmt.Println("got", n, "bytes.")
		buf = append(buf, tmp[:n]...)
	}
	return buf, true
}

func main() {

	url := "/download/"
	addr := "tortoisegit.org"
	port := 80
	data, flag := request_direct(url, addr, port, 3)
	if flag {
		fmt.Println(string(data))
	}
}
