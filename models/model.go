package models

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
