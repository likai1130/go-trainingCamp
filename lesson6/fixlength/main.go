package main

import (
	"log"
	"net"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	BYTES = 1024
)

//ClientTcpFixLength 客户端 fix length
func ClientTcpFixLength(conn net.Conn, sendMessage string) (int,error){
	sendByte := make([]byte, BYTES)
	temp := []byte(sendMessage)
	for j := 0; j < len(temp) && j < BYTES; j++ {
		sendByte[j] = temp[j]
	}
	return conn.Write(sendByte)
}

func client(){
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		log.Printf("Dial Error: %s (Client)\n", err)
		return
	}
	defer conn.Close()
	log.Printf("Connected to server. (remote address: %s, local address: %s) (Client)\n", conn.RemoteAddr(), conn.LocalAddr())
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		ClientTcpFixLength(conn, msg)
	}
}

func ServerTcpFixLength(conn net.Conn) {
	defer conn.Close()
	for {
		var buf = make([]byte, BYTES)
		_, err := conn.Read(buf)
		if err != nil {
			log.Printf("The connertion is closed by another side. (Server) \n")
			break
		}
		log.Printf("Received request: %s (Server)\n", string(buf))
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
		go ServerTcpFixLength(conn)
	}
}

func main() {
	go server()
	time.Sleep(500 * time.Millisecond)
	go client()
	time.Sleep(500 * time.Millisecond)
}
