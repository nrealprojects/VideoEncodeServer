package main

import (
	"path/filepath"
	Encode "projects/VideoEncodeServer/videoEncode"
)

var (
	// InputBasePath : InputBasePath
	InputBasePath = `C:\Users\nreal\Desktop\RecordRes\NetImages\input`
	// OutPutBasePath : OutPutBasePath
	OutPutBasePath = `C:\Users\nreal\Desktop\RecordRes\NetImages\output`
)

func main() {
	input := filepath.Join(InputBasePath, "275820104215")
	output := filepath.Join(OutPutBasePath, "275820104215")
	Encode.Do(input, output)
}
