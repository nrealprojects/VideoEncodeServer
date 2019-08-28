package videoencoder

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"projects/VideoEncodeServer/fileUtils"
	"strings"
)

var (
	// ToolPath : encode tool path
	ToolPath = `C:\Users\nreal\go\src\projects\VideoEncodeServer\encodeTools`
	// BlackList : BlackList
	BlackList = make(map[string]int, 0)
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
	totalline := 0
	for _, item := range files {
		_, exist := BlackList[item]
		if !exist {
			itempath := filepath.Join(input, item)
			line = fmt.Sprintf("file '%s'\n", itempath)
			buffer.WriteString(line)
			buffer.WriteString("duration 1\n")
			totalline++
		}
	}
	fmt.Println("Total line:", totalline)
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
	fileUtils.EnsureFolderExist(output)
	videoInput := filepath.Join(output, filename+"_input.txt")
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

func kickIllegalPicture(rgbinput, virtualinput string) {
	rgbfiles, err := fileUtils.GetAllFile(rgbinput)
	if err != nil {
		fmt.Println("Get all files err:", err)
		return
	}
	virtualfiles, err := fileUtils.GetAllFile(virtualinput)
	if err != nil {
		fmt.Println("Get all files err:", err)
		return
	}

	for _, rgb := range rgbfiles {
		result := false
		for _, virtual := range virtualfiles {
			if rgb == virtual {
				result = true
				break
			}
		}

		if !result {
			BlackList[rgb] = 0
			// 	fmt.Println("can not find rgb file:", rgb, " in virtual")
		}
	}

	for _, virtual := range virtualfiles {
		result := false
		for _, rgb := range rgbfiles {
			if rgb == virtual {
				result = true
				break
			}
		}
		if !result {
			BlackList[virtual] = 1
			// 	fmt.Println("can not find virtual file:", virtual, " in rgb")
		}
	}

	fmt.Println("Black list count:", len(BlackList))
}

// Do : Encode the input images to output
func Do(input string, output string) {
	rgbinput := filepath.Join(input, "RGB")
	virtualinput := filepath.Join(input, "Virtual")
	kickIllegalPicture(rgbinput, virtualinput)
	genVideoByImages(rgbinput, output)
	genVideoByImages(virtualinput, output)
}
