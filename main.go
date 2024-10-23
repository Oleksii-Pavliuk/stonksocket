package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

const port int = 8080

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error starting server", err)
	}
	defer listener.Close()

	log.Printf("Server started on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection", err)
		}
		go handleWebSocket(conn)
	}
}

func handleWebSocket(conn net.Conn) {
	reader := bufio.NewReader(conn)

	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println("Failed to read WebSocket upgrade request:", err)
		conn.Close()
		return
	}

	if strings.ToLower(req.Header.Get("Connection")) != "upgrade" || strings.ToLower(req.Header.Get("Upgrade")) != "websocket" {
		log.Println("Not upgade request:", err)
		conn.Close()
		return
	}

	webSocketKey := req.Header.Get("Sec-WebSocket-Key")
	if webSocketKey == "" {
		log.Println("Missing WebSocket key:", err)
		conn.Close()
		return
	}

	responseKey := computeAcceptKey(webSocketKey)
	response := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Sec-WebSocket-Accept: " + responseKey + "\r\n\r\n"
	conn.Write([]byte(response))

	log.Println("Client connected")

	for {
		message := make([]byte, 512)
		n, err := conn.Read(message)
		if err != nil {
			log.Println("Failed to read message", err)
			conn.Close()
			return
		}
		log.Printf("Recieved: %s\n", message[:n])

		conn.Write(encodeWebSocketMessage([]byte("Hello from server!")))
	}
}

func computeAcceptKey(webSocketKey string) string {
	websocketGUID := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(webSocketKey + websocketGUID))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func encodeWebSocketMessage(message []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x81)
	buf.WriteByte(byte(len(message)))
	buf.Write(message)
	return buf.Bytes()
}
