package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t' //将特殊的分隔符\t作为消息分隔符
)

/**
服务端读客户端数据
*/
func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

/**
处理连接
*/
func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		strReq, err := read(conn)
		if err != nil {
			if err == io.EOF {
				log.Printf("The connertion is closed by another side. (Server) \n")
			} else {
				log.Printf("Read Error: %s (Server)\n", err)
			}
			break
		}
		log.Printf("Received request: %s (Server)\n", strReq)
	}
}

func server() {
	var listener net.Listener
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		log.Printf("Listen Error :%s\n", err.Error())
		return
	}
	defer listener.Close()
	log.Printf("Got listener for the server.(local address: %s)\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept Error: %s \n", err.Error())
		}
		log.Printf("Established a connection with a client application. (remote address: %s)\n", conn.RemoteAddr())
		go handleConn(conn)
	}
}

/**
客户端把数据写给服务端
*/
func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}

/**
客户端
*/
func client() {
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		log.Printf("Dial Error: %s (Client)\n", err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to server. (remote address: %s, local address: %s) (Client)\n", conn.RemoteAddr(), conn.LocalAddr())

	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		write(conn, msg)
	}
}

func main() {
	go server()
	time.Sleep(500 * time.Millisecond)
	go client()
	time.Sleep(500 * time.Millisecond)
}
