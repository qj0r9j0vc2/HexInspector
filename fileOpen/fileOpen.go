package fileOpen

import (
	"HexInspector/util"
	"fmt"
	"os"
)

func FileOpen() {
	fileName := util.ReadFilePath()
	_, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Cannot open file: " + fileName)
		fmt.Println(err.Error())
	} else {
		fmt.Println("File open success!!: " + fileName)
	}
}
