package imageserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

var (
	// PackageSize : default package sieze
	PackageSize = 1300
	bufPool     sync.Pool
)

// Get : get a new packdata
func Get() PackData {
	v := bufPool.Get()
	var pack PackData
	if v == nil {
		//若不存在buf，创建新的
		pack.Init(PackageSize)
	} else {
		// 池里存在buf,v这里是interface{}，需要做类型转换
		pack = v.(PackData)
	}
	return pack
}

// Put : Put a new packdata
func Put(bf PackData) {
	bufPool.Put(bf)
}

// PackData :
type PackData struct {
	DataSize   int    //數據大小
	CurSize    int    //當前大小
	IsLastPack int    //是否是最後一片
	Data       []byte //帧数据
}

// Init : init
func (p *PackData) Init(datasize int) {
	p.Data = make([]byte, datasize)
}

// ReadOnePack : read one pack
func (p *PackData) ReadOnePack(data []byte) {
	// fmt.Println(" start to read one pack data len:", len(data))
	p.DataSize = BytesToInt(data[0:4])
	p.CurSize = BytesToInt(data[4:8])
	p.IsLastPack = BytesToInt(data[8:12])
	p.Data = data[12:(p.CurSize + 12)]
	// fmt.Println("data len:", len(data), "this pack is :", p.ToString())
}

// ToString : ToString
func (p *PackData) ToString() string {
	return fmt.Sprintf("totalsize: %d cursize: %d datasize:%d is last pak:%d", p.DataSize, p.CurSize, len(p.Data), p.IsLastPack)
}

// BytesToInt : BytesToInt
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

// ToBytes : to bytes
func (p *PackData) ToBytes() []byte {
	var buffer bytes.Buffer
	buffer.Write(IntToBytes(p.DataSize))
	buffer.Write(IntToBytes(p.CurSize))
	buffer.Write(IntToBytes(p.IsLastPack))
	buffer.Write(p.Data)
	return buffer.Bytes()
}

// IntToBytes : int to bytes
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

// FrameData :
type FrameData struct {
	Imagetype int    //數據大小
	TimeStamp uint64 //當前大小
	Data      []byte //帧数据
}

// ToBytes : to bytes
func (p *FrameData) ToBytes() []byte {
	var buffer bytes.Buffer
	buffer.Write(IntToBytes(p.Imagetype))
	buffer.Write(UInt64ToBytes(p.TimeStamp))
	buffer.Write(p.Data)
	return buffer.Bytes()
}

// ToString : ToString
func (p *FrameData) ToString() string {
	return fmt.Sprintf("Imagetype: %d TimeStamp: %d Data len:%d ", p.Imagetype, p.TimeStamp, len(p.Data))
}

// BytesToFrame : Bytes To Frame
func BytesToFrame(data []byte) FrameData {
	var frame FrameData
	frame.Imagetype = BytesToInt(data[0:4])
	frame.TimeStamp = uint64(binary.LittleEndian.Uint64(data[4:12]))
	frame.Data = data[12:(len(data))]
	return frame
}

// UInt64ToBytes : uint64 to bytes
func UInt64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}
