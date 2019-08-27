package fileUtils

import (
	"testing"
)

// func TestFileReader(t *testing.T) {
// 	len := genInputfile()
// 	fmt.Println(len)
// }

// func TestGeninputFile(t *testing.T) {
// 	poselist := genRGBInputfile()
// 	genVirtualInputfile(poselist)
// }

func TestCreateDirectory(t *testing.T) {
	err := EnsureFolderExist(`C:\Users\nreal\Desktop\RecordRes\NetImages\input\` + "23234234234234")
	if err != nil {
		t.Error(err)
	}
}

// func TestFileReaderLine(t *testing.T) {
// 	path := `C:\\Users\\nreal\\Desktop\\RecordRes\\propagator_pose.dat`
// 	result := ReadFileLine(path)
// 	if result != nil {
// 		fmt.Println("success :", len(result))
// 	} else {
// 		t.Error("read error")
// 	}
// }

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// func TestJson(t *testing.T) {
// 	s := Student{
// 		Name: "aze",
// 		Age:  12,
// 	}
// 	data, _ := json.Marshal(s)
// 	fmt.Println(string(data))

// 	file := []byte(`{
// 		"name": "attila@attilaolah.eu",
// 		"age": 12,
// 	  }`)
// 	student := &Student{}
// 	json.Unmarshal(file, student)
// 	fmt.Println(*student)
// }
