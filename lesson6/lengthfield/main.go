package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

/**
LengthFieldBasedFrameDecoder 在数据包中添加长度字段
*/
const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
)

// encode 将消息编码
func encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// decode 将消息解码
func decode(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4) // 读取前4个字节的数据
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// Buffered返回缓冲中现有的可读取的字节数。
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}
	// 读取真正的消息数据
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := decode(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("read from client failed, err: %v\n", err)
			break
		}
		log.Printf("收到client发来的数据：%s\n", msg)
	}
}

//server 模拟服务端请求
func server() {
	listen, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		log.Printf("listen failed, err:%v \n", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept() //建立连接
		if err != nil {
			log.Printf("accept failed, err: %v\n", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}

// client 模拟客户端请求
func client() {
	conn, err := net.Dial(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		log.Printf("dial failed, err: %v \n", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		data, err := encode(msg)
		if err != nil {
			log.Printf("encode msg failed, err: %v \n", err)
			return
		}
		conn.Write(data)
	}
}

func main() {
	go server()
	time.Sleep(500 * time.Millisecond)
	go client()
	time.Sleep(500 * time.Millisecond)
}
