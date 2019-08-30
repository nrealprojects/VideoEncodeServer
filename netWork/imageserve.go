package imageserver

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"projects/VideoEncodeServer/fileUtils"
	"time"
)

var (
	recChanel = make(chan []byte, 1024*1024)
)

// Start : 模拟server端
func Start(encodech chan string, outputpath string) {
	fmt.Println("Start image server....")
	tcpServer, _ := net.ResolveTCPAddr("tcp4", ":6000")
	listener, _ := net.ListenTCP("tcp", tcpServer)

	go reciveMsg(recChanel)
	go encode(CompletPackageCh, encodech, outputpath)

	for {
		fmt.Println("start a connect!")
		//当有新的客户端请求来的时候，拿到与客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//处理逻辑
		go handle(conn, recChanel)
	}
}

func handle(conn net.Conn, c chan []byte) {
	// defer conn.Close()
	//读取客户端传送的消息
	fmt.Println("a client has connected!!")
	go func() {
		data := make([]byte, 1024*1024)
		for {
			i, err := conn.Read(data)
			if err != nil {
				fmt.Println("读取客户端数据错误:", err.Error())
				break
			}
			copydata := make([]byte, i)
			copy(copydata, data[0:i])
			c <- copydata
		}
	}()
}

// encode : encode the receive data
func encode(c chan FrameData, preEncoded chan string, outputpath string) {
	var clearZero uint64
	var newPath string
	var newfolder string
	clearZero = 0
	for {
		select {
		case frame := <-c:
			// if current frame is a end signal
			// then add current folder to preencode listadb
			if frame.IsEndSignal() {
				if newfolder != "" {
					fmt.Println("--------rec a end signal-------")
					preEncoded <- newfolder
					newfolder = ""
					clearZero = 0
				}
			} else {
				folder := "RGB"
				if frame.Imagetype == 0 {
					folder = "RGB"
				} else if frame.Imagetype == 1 {
					folder = "Virtual"
				} else if frame.Imagetype == 2 {
					folder = "Unknow"
				}
				if clearZero == 0 {
					// create a new folder
					newfolder = createNewFolder()
					newPath = filepath.Join(outputpath, fmt.Sprintf("/%s", newfolder))
					fileUtils.EnsureFolderExist(newPath)
					fileUtils.EnsureFolderExist(filepath.Join(newPath, "RGB"))
					fileUtils.EnsureFolderExist(filepath.Join(newPath, "Virtual"))
					fileUtils.EnsureFolderExist(filepath.Join(newPath, "Unknow"))
				}

				imageSavepath := filepath.Join(newPath, fmt.Sprintf("%s\\%d.jpg", folder, frame.TimeStamp))
				err := ioutil.WriteFile(imageSavepath, frame.Data, 0644)
				if err != nil {
					fmt.Println("Write image err:", err)
				}
				clearZero = frame.TimeStamp
			}
		case <-time.After(8 * time.Second):
			if clearZero != 0 {
				fmt.Println("-------clear-------")
				clearZero = 0

				// add the new folder to encode thread
				if newfolder != "" {
					preEncoded <- newfolder
					newfolder = ""
				}
			}
		}
	}
}

func createNewFolder() string {
	return time.Now().Format("2006-01-02-15-04-05")
}
