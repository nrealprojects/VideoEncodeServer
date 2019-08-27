package fileUtils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// IsPathExists : return whether a file is exist
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// EnsureFileNotExist : Ensure File is Not Exist
func EnsureFileNotExist(path string) error {
	result, err := IsPathExists(path)
	if !result {
		os.Remove(path)
	}
	return err
}

// EnsureFolderExist : Ensure folder is Exist
func EnsureFolderExist(path string) error {
	result, err := IsPathExists(path)
	if !result {
		err = os.MkdirAll(path, os.ModePerm)
	}
	fmt.Println("create directory :", path, err)
	return err
}

// ReadFile : read a file to bytes
func ReadFile(path string) []byte {
	//打开文件
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("err = ", err)
		return nil
	}

	//关闭文件
	defer f.Close()
	buf := make([]byte, 1024*2000) //2k大小

	//n代表从文件读取内容的长度
	n, err1 := f.Read(buf)
	if err1 != nil && err1 != io.EOF { //文件出错，同时没有到结尾
		fmt.Println("err1 = ", err1)
		return nil
	}
	// fmt.Println("len of file:", n)
	// fmt.Println(string(buf[:n]))
	return buf[:n]
}

// ReadFileLine : Read a file per line
func ReadFileLine(path string) []string {
	//打开文件
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("err = ", err)
		return nil
	}

	//关闭文件
	defer f.Close()

	result := make([]string, 0)
	//新建一个缓冲区，把内容先放在缓冲区
	r := bufio.NewReader(f)
	for {
		//遇到'\n'结束读取, 但是'\n'也读取进入
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF { //文件已经结束
				break
			}
			fmt.Println("err = ", err)
		}
		result = append(result, string(buf))
	}

	return result
}

// GetAllFile : get all files of folder
func GetAllFile(folder string) ([]string, error) {
	files := make([]string, 0)
	rd, err := ioutil.ReadDir(folder)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", folder+"\\"+fi.Name())
			GetAllFile(folder + fi.Name() + "\\")
		} else {
			files = append(files, fi.Name())
		}
	}
	return files, err
}
