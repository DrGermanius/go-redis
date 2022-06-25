package main

import (
	"fmt"
	"net"
	"os"
)

var cache *Cache

func main() {
	cache = NewCache()
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Printf("Failed to bind to port 6379: %s\n", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			continue
		}
		go serve(conn)
	}
}

func serve(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)

		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("failed to read from socket: %s\n", err)
			return
		}

		var res []byte
		switch {
		case isRESPArray(buffer):
			res = processRESPArray(buffer)
		default:
			res = encodeToRESPSimpleString("PONG")
		}

		_, err = conn.Write(res)
		if err != nil {
			fmt.Printf("failed to read from socket: %s\n", err)
			return // EOF not a problem
		}
	}
}
