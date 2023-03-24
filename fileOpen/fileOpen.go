package fileOpen

import (
	"HexInspector/util"
	"fmt"
	"os"
)

func FileOpen() {
	fileName := util.ReadFilePath()
	_, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Cannot Open file: " + err.Error())
		return
	}
	util.Clear()

}
