package main

import (
	"bytes"
	"fmt"
)

var (
	dataSizePerFrame = 1024
	// CompletPackageCh ： CompletPackageCh
	CompletPackageCh = make(chan []byte, 1024*1024)
)

// ReciveMsg : 接收服务器数据
func ReciveMsg(ch chan []byte) {
	packlist := make([]PackData, 0)
	for {
		select {
		case data := <-ch:
			count := len(data) / dataSizePerFrame
			if len(data)%dataSizePerFrame != 0 {
				count++
			}
			// fmt.Println("Read a pack len:", len(data), " count is:", count)
			for i := 0; i < count; i++ {
				var pack PackData
				if i != count-1 {
					pack.ReadOnePack(data[i*dataSizePerFrame : (i+1)*dataSizePerFrame])
				} else {
					pack.ReadOnePack(data[i*dataSizePerFrame : len(data)])
				}

				// fmt.Println("Read a pack success :", pack.ToString())
				if pack.IsLastPack == 1 {
					// join the packlist to a full package and clear the packlist
					packlist = append(packlist, pack)
					fullpack := JoinPackData(packlist)
					if fullpack != nil {
						CompletPackageCh <- fullpack
						// fmt.Println("Success get a full package.")
					} else {
						fmt.Println("Faild get a full package.")
					}
					packlist = packlist[:0]
					// fmt.Println("clear pack list count:", len(packlist))
				} else {
					packlist = append(packlist, pack)
				}
			}
		}
	}
}

// JoinPackData ：joint pack list
func JoinPackData(list []PackData) []byte {
	// fmt.Println("JoinPackData list count:", len(list))
	framesize := 0
	var frameBuffer []byte
	var buffer bytes.Buffer
	for i := 0; i < len(list); i++ {
		data := list[i]
		framesize = data.DataSize
		buffer.Write(data.Data)
	}
	frameBuffer = buffer.Bytes()

	// fmt.Println("join the packages frame size:", framesize, "real len:", len(frameBuffer))
	if framesize != len(frameBuffer) {
		return nil
	}
	return frameBuffer
}
