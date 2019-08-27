package models

import "fmt"

// PathInfo :
type PathInfo struct {
	// CameraMatrix    string     `json:"cameraMatrix"`
	TranPoses []PoseItem `json:"tranPoses"`
	// CameraLocalPose string     `json:"cameraLocalPose"`
}

// PoseItem :
type PoseItem struct {
	TimeNanos uint64 `json:"timeNanos"`
	// Items     string `json:"items"`
}

// SlamPose :
type SlamPose struct {
	TimeNanos uint64
}

// PackData :
type PackData struct {
	datasize  int    //數據大小
	cursize   int    //當前大小
	islastpak int    //是否是最後一片
	data      []byte //帧数据
}

// Init : init
func (p *PackData) Init(datasize int) {
	p.data = make([]byte, datasize)
}

// ToString : ToString
func (p *PackData) ToString() string {
	return fmt.Sprintf("datasize: %d cursize: %d is last pak:%d ", p.datasize, p.cursize, p.islastpak)
}
