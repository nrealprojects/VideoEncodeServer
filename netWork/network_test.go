package imageserver

import (
	"fmt"
	"testing"
	"time"
)

// func TestPool(t *testing.T) {
// 	return
// 	ch := make(chan PackData, 100)
// 	go func() {
// 		for i := 0; i < 10; i++ {
// 			pack := Get()
// 			pack.Init(1300)
// 			ch <- pack
// 		}
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case pack := <-ch:
// 				fmt.Println(pack.ToString())
// 				Put(pack)
// 			case <-time.After(1 * time.Second):
// 				return
// 			}
// 		}
// 	}()

// 	time.Sleep(2 * time.Second)
// }

// func TestByteConvert(t *testing.T) {
// 	intbytes := IntToBytes(12)
// 	fmt.Println("intbytes count:", len(intbytes))

// 	var pack PackData
// 	pack.DataSize = 1024 * 10
// 	pack.CurSize = 1300
// 	pack.Data = make([]byte, 1024)
// 	pack.IsLastPack = 1
// 	packbytes := pack.ToBytes()
// 	fmt.Println("packbytes count:", len(packbytes))
// }

func TestPackageSender(t *testing.T) {
	var frame FrameData
	frame.Imagetype = 2
	frame.TimeStamp = 12
	frame.Data = make([]byte, 10240)

	framebyte := frame.ToBytes()
	fmt.Println("frame bytes count:", len(framebyte))

	datach := make(chan []byte, 100)
	go func(data []byte) {
		dataSize := dataSizePerFrame - 12
		count := len(data) / dataSize
		if len(data)%dataSize != 0 {
			count++
		}
		for i := 0; i < count; i++ {
			pack := new(PackData)
			pack.DataSize = len(data)
			if i != count-1 {
				pack.CurSize = dataSizePerFrame
				pack.IsLastPack = 0
				pack.Data = make([]byte, dataSize)
				copy(pack.Data, data[i*dataSize:(i+1)*dataSize])
			} else {
				pack.CurSize = len(data) % dataSize
				pack.IsLastPack = 1
				pack.Data = make([]byte, len(data)%dataSize)
				copy(pack.Data, data[i*dataSize:(i*dataSize+len(data)%dataSize)])
			}
			packbytes := pack.ToBytes()
			fmt.Println("test pack is:", pack.ToString())
			datach <- packbytes
		}
	}(framebyte)

	go ReciveMsg(datach)

	go func() {
		for {
			select {
			case data := <-CompletPackageCh:
				frame := BytesToFrame(data)
				fmt.Println(frame.ToString())
			case <-time.After(1 * time.Second):
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)
}
