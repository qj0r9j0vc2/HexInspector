package util

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	//Colors
	FgHiGreenUnder = color.New(color.FgHiGreen).Add(color.Underline)
	FgHiCyanUnder  = color.New(color.FgHiCyan).Add(color.Underline)
	FgHiCyan       = color.New(color.FgHiCyan)
	FgHiRed        = color.New(color.FgRed)
)

func Clear() {
	c := exec.Command("Clear")
	c.Stdout = os.Stdout
	c.Run()
}

func InputStr() string {
	reader := bufio.NewReader(os.Stdin)

	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("입력된 데이터가 올바르지 않습니다.: %s \n", err.Error())
	}
	str = strings.TrimSpace(str)

	return str
}

func ReadFilePath() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter file path : ")
	strName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fileName := strings.TrimSpace(strName)
	return fileName
}

func StringToHex(str string) []int {

	var intArr []int
	var sumArr []int
	for i := 0; i < len(str); i++ {
		byte := int(str[i]) - 48
		if byte >= 17 {
			if byte >= 49 {
				intArr = append(intArr, byte-39)
			} else {
				intArr = append(intArr, byte-7)
			}
		} else if byte < 10 {
			intArr = append(intArr, byte)
		} else {
			intArr = append(intArr, 0)
		}
		if i%2 == 1 {
			sumArr = append(sumArr, intArr[i-1]*16+intArr[i])
		}
	}

	return sumArr
}
