package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"projects/encodeImages/fileUtils"
	"projects/encodeImages/models"
	"strconv"
	"strings"
)

var (
	// ImageBasePath : image base path
	ImageBasePath        string = `C:\Users\nreal\Desktop\RecordRes\RecordImages\RGB`
	VirtualImageBasePath string = `D:\WorkSpace\Projects\GitLab\NRSDKForUnity\RecordImages\Virtual`
	// CurrntPath : image base path
	CurrntPath string = `C:\Users\nreal\go\src\projects\encodeImages`
	// RgbTimeStamFile : image base path
	RgbTimeStamFile string = "C:\\Users\\nreal\\Desktop\\RecordRes\\pathinfo.txt"
	// SlamTimeStamFile : image base path
	SlamTimeStamFile string = "C:\\Users\\nreal\\Desktop\\RecordRes\\propagator_pose.dat"
	// videoInputFile : image base path
	rgbVideoInputFile string = `C:\Users\nreal\Desktop\RecordRes\rgbinput.txt`
	// videoInputFile : image base path
	virtualVideoInputFile string = `C:\Users\nreal\Desktop\RecordRes\virtualinput.txt`
	// outPutFile : image base path
	outPutFile string = `C:\Users\nreal\Desktop\RecordRes`
)

func genRGBInputfile() []models.PoseItem {
	file, _ := ioutil.ReadFile(RgbTimeStamFile)
	fmt.Println("file len:", len(file))
	pathInfo := &models.PathInfo{}
	err := json.Unmarshal(file, pathInfo)
	if err != nil {
		fmt.Println("some thing error:", err)
	}
	fmt.Println("read success:", len(pathInfo.TranPoses))

	var buffer bytes.Buffer
	var line string
	for _, item := range pathInfo.TranPoses {
		itempath := filepath.Join(ImageBasePath, fmt.Sprintf("%d.png", item.TimeNanos))
		line = fmt.Sprintf("file '%s'\n", itempath)
		buffer.WriteString(line)
		buffer.WriteString("duration 1\n")
	}
	buffer.WriteString(line)
	// defer fmt.Println(buffer.String())
	if ioutil.WriteFile(rgbVideoInputFile, buffer.Bytes(), 0644) == nil {
		fmt.Println("Write to file success!")
	}
	return pathInfo.TranPoses
}

func genRGBVideo() []models.PoseItem {
	//gen rgbinput fle by pathinfo(RGB timestamp)
	rgbposes := genRGBInputfile()

	fmt.Println("start to run encode bat")
	encodbat := filepath.Join(CurrntPath, "encode.bat")
	flipbat := filepath.Join(CurrntPath, "flip.bat")
	fmt.Println(encodbat)
	if result, _ := fileUtils.IsPathExists(encodbat); !result {
		fmt.Println("bat file is not exist!")
	}

	outputpath := filepath.Join(outPutFile, "output.mp4")
	outputpath2 := filepath.Join(outPutFile, "output_flip.mp4")
	//delet old output.mp4
	fileUtils.EnsureFileNotExist(outputpath)
	fileUtils.EnsureFileNotExist(outputpath2)

	// generate Rgb video
	out, err := exec.Command(encodbat, rgbVideoInputFile, outputpath).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))

	out1, err1 := exec.Command(flipbat, outputpath, outputpath2).CombinedOutput()
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(string(out1))

	return rgbposes
}

func genVirtualInputfile(poslist []models.PoseItem) {
	file := fileUtils.ReadFileLine(SlamTimeStamFile)
	if file == nil {
		fmt.Println("read file error,quit!")
		os.Exit(-1)
	}
	fmt.Println("file len:", len(file))

	virtualimages := make([]models.SlamPose, 0)
	for _, item := range file {
		lines := strings.SplitAfter(item, ",")
		time := strings.Trim(lines[1], " ")
		time = strings.Trim(time, ",")

		timestamp, err := strconv.ParseUint(time, 10, 64)
		if err != nil {
			fmt.Println(err)
			break
		}
		var pose models.SlamPose
		pose.TimeNanos = timestamp
		virtualimages = append(virtualimages, pose)
	}

	virtualIndex := make([]uint64, 0)
	for _, item := range poslist {
		tempindex := BinarySearch(virtualimages, 0, len(virtualimages), item.TimeNanos)
		virtualIndex = append(virtualIndex, virtualimages[tempindex].TimeNanos)
	}

	var buffer bytes.Buffer
	var line string
	for _, item := range virtualIndex {
		itempath := filepath.Join(VirtualImageBasePath, fmt.Sprintf("%d.png", item))
		if exist, _ := fileUtils.IsPathExists(itempath); exist {
			line = fmt.Sprintf("file '%s'\n", itempath)
			buffer.WriteString(line)
			buffer.WriteString("duration 1\n")
		} else {
			fmt.Println("file is not exist :", itempath)
		}
	}
	buffer.WriteString(line)
	fileUtils.EnsureFileNotExist(virtualVideoInputFile)

	if ioutil.WriteFile(virtualVideoInputFile, buffer.Bytes(), 0644) == nil {
		fmt.Println("Write to virtualVideoInputFile success!")
	}
}

func genVirtualVideo(poslist []models.PoseItem) {
	//gen virtualinput fle by propagator_pose(virtual timestamp)
	genVirtualInputfile(poslist)

	fmt.Println("start to run encode bat")
	encodbat := filepath.Join(CurrntPath, "encode.bat")
	fmt.Println(encodbat)
	if result, _ := fileUtils.IsPathExists(encodbat); !result {
		fmt.Println("bat file is not exist!")
	}

	outputpath := filepath.Join(outPutFile, "virtual_output.mp4")
	//delet old virtual_output.mp4
	fileUtils.EnsureFileNotExist(outputpath)

	// generate virtual video
	out, err := exec.Command(encodbat, virtualVideoInputFile, outputpath).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func main() {
	poselist := genRGBVideo()
	genVirtualVideo(poselist)
}
