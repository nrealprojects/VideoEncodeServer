package videoencoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"projects/encodeImages/fileUtils"
	"strings"
)

var (
	// ImageBasePath : image base path
	ImageBasePath = `C:\Users\nreal\Desktop\RecordRes\RecordImages\RGB`
	// ToolPath : encode tool path
	ToolPath = `C:\Users\nreal\go\src\projects\encodeImages\encodeTools`
	// outPutFile : image base path
	// outPutFile       = `C:\Users\nreal\Desktop\RecordRes`
	outPutImagesPath = `C:\Users\nreal\Desktop\RecordRes\NetImages\`
)

func genVideoByImages(input string, output string) {
	//outPutImagesPath + "RGB"
	//gen video input file
	files, err := fileUtils.GetAllFile(input)
	if err != nil {
		fmt.Println("Get all files err:", err)
		return
	}

	var buffer bytes.Buffer
	var line string
	for _, item := range files {
		itempath := filepath.Join(input, item)
		line = fmt.Sprintf("file '%s'\n", itempath)
		buffer.WriteString(line)
		buffer.WriteString("duration 1\n")
	}
	buffer.WriteString(line)

	filename := "RGB"
	if strings.Contains(input, "RGB") {
		filename = "RGB"
	} else if strings.Contains(input, "Virtual") {
		filename = "Virtual"
	} else if strings.Contains(input, "Unknow") {
		filename = "Unknow"
	}

	// defer fmt.Println(buffer.String())
	videoInput := output + filename + "_input.txt"
	fmt.Println(videoInput)
	if ioutil.WriteFile(videoInput, buffer.Bytes(), 0644) == nil {
		fmt.Println("Write to file success!")
	}

	// gen video by input file
	flip := false
	if strings.Contains(filename, "RGB") {
		flip = true
	}

	genVideoByInput(videoInput, output, flip)
}

func genVideoByInput(input string, output string, flip bool) {
	fmt.Println("start to run encode bat")
	encodbat := filepath.Join(ToolPath, "encode.bat")

	fmt.Println(encodbat)
	if result, _ := fileUtils.IsPathExists(encodbat); !result {
		fmt.Println("bat file is not exist!")
	}

	var videoFile string
	var videoFlipFile string
	if flip {
		videoFile = filepath.Join(output, "RGB.mp4")
		videoFlipFile = filepath.Join(output, "RGB_flip.mp4")
	} else {
		videoFile = filepath.Join(output, "Virtual.mp4")
	}

	//delet old output.mp4
	fileUtils.EnsureFileNotExist(output)
	// generate Rgb video
	_, err := exec.Command(encodbat, input, videoFile).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	if flip {
		flipbat := filepath.Join(ToolPath, "flip.bat")
		fileUtils.EnsureFileNotExist(videoFlipFile)
		_, err1 := exec.Command(flipbat, videoFile, videoFlipFile).CombinedOutput()
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}

// Do : Encode the input images to output
func Do(input string, output string) {
	rgbinput := filepath.Join(input, "RGB")
	virtualinput := filepath.Join(input, "Virtual")
	genVideoByImages(rgbinput, output)
	genVideoByImages(virtualinput, output)
}
