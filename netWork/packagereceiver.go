package main

import (
	"bytes"
)

var (
	dataSizePerFrame = 1600
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
			for i := 0; i < count; i++ {
				var pack PackData
				if (i+1)*dataSizePerFrame <= len(data) {
					pack.ReadOnePack(data[i*dataSizePerFrame : (i+1)*dataSizePerFrame])
				} else {
					pack.ReadOnePack(data[i*dataSizePerFrame : len(data)-i*dataSizePerFrame])
				}
				if pack.IsLastPack == 1 {
					// join the packlist to a full package and clear the packlist
					packlist := append(packlist, pack)
					CompletPackageCh <- JoinPackData(packlist)
					packlist = make([]PackData, 0)
				} else {
					packlist = append(packlist, pack)
				}
			}
		}
	}
}

// JoinPackData ：joint pack list
func JoinPackData(list []PackData) []byte {
	framesize := 0
	var frameBuffer []byte
	var buffer bytes.Buffer
	for i := 0; i < len(list); i++ {
		data := list[i]
		framesize = data.DataSize
		buffer.Write(data.Data)
	}
	frameBuffer = buffer.Bytes()
	if framesize != len(frameBuffer) {
		return nil
	}
	return frameBuffer
}
