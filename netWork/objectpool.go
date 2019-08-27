package main

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
	if v == nil {
		//若不存在buf，创建新的
		var pack PackData
		pack.Init(PackageSize)
		return pack
	} else {
		// 池里存在buf,v这里是interface{}，需要做类型转换
		return v.(PackData)
	}
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
	var pack PackData
	datasizebuf := bytes.NewBuffer(data[0:4])
	binary.Read(datasizebuf, binary.LittleEndian, &pack.DataSize)
	cursizebuf := bytes.NewBuffer(data[4:8])
	binary.Read(cursizebuf, binary.LittleEndian, &pack.CurSize)
	lastpackbuf := bytes.NewBuffer(data[8:12])
	binary.Read(lastpackbuf, binary.LittleEndian, &pack.IsLastPack)
	pack.Data = data[12:(len(data) - 12)]
}

// ToString : ToString
func (p *PackData) ToString() string {
	return fmt.Sprintf("datasize: %d cursize: %d is last pak:%d ", p.Data, len(p.Data), p.IsLastPack)
}
