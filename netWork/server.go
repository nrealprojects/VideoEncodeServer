package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"projects/encodeImages/fileUtils"
	"time"
)

var (
	recChanel = make(chan []byte, 1024*1024)
	// ResBasePath : Resource Base Path
	ResBasePath = `C:\Users\nreal\Desktop\RecordRes\NetImages\input\`
)

//模拟server端
func main() {
	tcpServer, _ := net.ResolveTCPAddr("tcp4", ":6000")
	listener, _ := net.ListenTCP("tcp", tcpServer)

	go Encode(recChanel)

	for {
		fmt.Println("start a connect!")
		//当有新的客户端请求来的时候，拿到与客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//处理逻辑
		go handle(conn, recChanel)
	}
}

func handle(conn net.Conn, c chan []byte) {
	// defer conn.Close()
	//读取客户端传送的消息
	fmt.Println("a client has connected!!")
	go func() {
		data := make([]byte, 512*1024)
		for {
			i, err := conn.Read(data)
			if err != nil {
				fmt.Println("读取客户端数据错误:", err.Error())
				break
			}
			copydata := make([]byte, i)
			copy(copydata, data[0:i])
			c <- copydata
		}
	}()
}

// Encode : encode the receive data
func Encode(c chan []byte) {
	var clearZero uint64
	var newPath string
	clearZero = 0
	for {
		select {
		case data := <-c:
			var imagetype int32
			binbuf := bytes.NewBuffer(data[0:4])
			binary.Read(binbuf, binary.LittleEndian, &imagetype)
			timestamp := uint64(binary.LittleEndian.Uint64(data[4:12]))
			folder := "RGB"
			if imagetype == 0 {
				folder = "RGB"
			} else if imagetype == 1 {
				folder = "Virtual"
			} else if imagetype == 2 {
				folder = "Unknow"
			}
			if clearZero == 0 {
				newPath = ResBasePath + fmt.Sprintf("%d", timestamp)
				fileUtils.EnsureFolderExist(newPath)
				fileUtils.EnsureFolderExist(filepath.Join(newPath, "RGB"))
				fileUtils.EnsureFolderExist(filepath.Join(newPath, "Virtual"))
				fileUtils.EnsureFolderExist(filepath.Join(newPath, "Unknow"))
			}
			imageSavepath := filepath.Join(newPath, fmt.Sprintf("%s\\%d.jpg", folder, timestamp))
			ioutil.WriteFile(imageSavepath, data[12:len(data)], 0644)
			fmt.Println("encode a image:", folder, " ", timestamp)
			// if writeErr == nil {
			// 	fmt.Println("Save a client screen shot image success!", imageSavepath)
			// } else {
			// 	fmt.Println("Save a client screen shot image failed :", writeErr)
			// }
			clearZero = timestamp
		case <-time.After(8 * time.Second):
			fmt.Println("-------clear-------")
			clearZero = 0
		}

	}
}
