package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

var outputPath = "./output"

//处理路径
func GetFilePath(fileName string) (fullname string, err error) {
	dir, filename := filepath.Split(fileName)
	//获取程序运行时候的相对路径
	//tmpDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	return "", err
	// }
	if dir != "" {
		dir = outputPath + dir
	} else {
		dir = outputPath
	}

	fullname = filepath.Join(dir, filename)
	//判断文件夹是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777) //0777也可以os.ModePerm
	}
	fmt.Printf("GetFilePath fullname %s", fullname)
	return
}
