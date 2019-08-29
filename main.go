package main

import (
	"log"
	"os"
	"path/filepath"
	"projects/VideoEncodeServer/fileUtils"
	Server "projects/VideoEncodeServer/netWork"
	Encode "projects/VideoEncodeServer/videoEncode"
	"time"
)

var (
	inputCh = make(chan string, 100)
)

func main() {
	curdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
		return
	}
	input := filepath.Join(curdir, "Generate/imagesOutput")
	output := filepath.Join(curdir, "Generate/videoOutput")
	fileUtils.EnsureFolderExist(input)
	fileUtils.EnsureFolderExist(output)

	// start server
	go Server.Start(inputCh, input)
	// start encode thread
	go Encode.Start(inputCh, input, output)

	for {
		time.Sleep(time.Microsecond)
	}
}
