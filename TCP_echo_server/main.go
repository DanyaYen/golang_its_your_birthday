package main

import (
	"net"
	"fmt"
	"io"
)

func handleConnection(conn net.Conn, messages chan string) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {

		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Read error", err)
			return
		}
		messages <- string(buf[:n])

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("Write error", err)
			return
		}
	}
	

}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	messages := make(chan string)
	go func() {
		for msg := range messages {
			fmt.Println("Received from client:", msg)
		}
	}()
	
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	
	defer listener.Close()
	for { 
		conn, err := listener.Accept() 
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, messages)
	}	
	
}
