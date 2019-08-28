package main

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	ch := make(chan PackData, 100)
	go func() {
		for i := 0; i < 10; i++ {
			pack := Get()
			pack.Init(1300)
			ch <- pack
		}
	}()

	go func() {
		for {
			select {
			case pack := <-ch:
				fmt.Println(pack.ToString())
				Put(pack)
			case <-time.After(1 * time.Second):
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)
}

func TestPackageSender(t *testing.T) {
	var frame FrameData
	frame.Imagetype = 2
	frame.TimeStamp = 123123123123
	frame.Data = make([]byte, 1024*1024)

	framebyte := frame.ToBytes()
	fmt.Println("frame bytes count:", len(framebyte))

	datach := make(chan []byte, 100)
	go func(data []byte) {
		count := len(data) / dataSizePerFrame
		if len(data)%dataSizePerFrame != 0 {
			count++
		}
		for i := 0; i < count; i++ {
			var pack PackData
			pack.DataSize = len(data)
			if i != count-1 {
				pack.CurSize = dataSizePerFrame
				pack.IsLastPack = 0
				pack.Data = data[i*dataSizePerFrame : (i+1)*dataSizePerFrame]
			} else {
				pack.CurSize = len(data) % dataSizePerFrame
				pack.IsLastPack = 1
				pack.Data = data[i*dataSizePerFrame : (i*dataSizePerFrame + len(data)%dataSizePerFrame)]
			}

			datach <- pack.ToBytes()
		}
	}(framebyte)
	go ReciveMsg(datach)

	go func() {
		select {
		case data := <-CompletPackageCh:
			frame := BytesToFrame(data)
			fmt.Println(frame.ToString())
		case <-time.After(1 * time.Second):
			break
		}
	}()

	time.Sleep(2 * time.Second)
}
