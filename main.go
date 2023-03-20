package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var content []byte
var sectorSize = int64(512)
var fgHiGreenUnder = color.New(color.FgHiGreen).Add(color.Underline)
var fgHiCyanUnder = color.New(color.FgHiCyan).Add(color.Underline)
var fgHiCyan = color.New(color.FgHiCyan)
var fgHiRed = color.New(color.FgRed)
var contentSize int64

type foundHexArray struct {
	startIdx int64
	endIdx   int64
}

func main() {
	menu()
}

func menu() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[1;1H")
	for {
		menu := []string{"File Open", "Show Sector Information", "Show Partition Information", "Exit"}
		prompt := promptui.Select{Label: "Select Menu", Items: menu}
		i, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		switch i {
		case 0:
			fileOpen()
			break
		case 1:
			showSectorInformation()
			break
		case 2:
			showPartitionInformation()
			break
		case 3:
			println("Exit..!!")
			return
		}
	}
}

func inputStr() string {
	reader := bufio.NewReader(os.Stdin)

	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("입력된 데이터가 올바르지 않습니다.: %s \n", err.Error())
	}
	str = strings.TrimSpace(str)

	return str
}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func readFilePath() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter file path : ")
	strName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fileName := strings.TrimSpace(strName)
	return fileName
}

func fileOpen() {
	fileName := readFilePath()
	_, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Cannot open file: " + fileName)
	} else {
		fmt.Println("File open success!!: " + fileName)
	}
}

func stringToHex(str string) []int {

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

func findSector() {
	lastStartIdx := int64(0)
	lastEndIdx := sectorSize
	maxSectorSize := contentSize / sectorSize

	clear()
	printBlock(content, 0, int64(sectorSize), []foundHexArray{})

	for {
		fgHiGreenUnder.Println("Enter sector idx")
		fmt.Printf("(if you want to search hex value, enter hex value after 'f')\n")
		str := inputStr()
		idx, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			if str == "f" || str == "F" {
				var str string
				str = inputStr()
				str = strings.ReplaceAll(str, " ", "")

				foundHexArray := findHexArray(content, lastStartIdx, lastEndIdx, stringToHex(str))

				clear()
				printBlock(content, lastStartIdx, lastEndIdx, foundHexArray)
			}
		} else {
			if idx > maxSectorSize {
				clear()

				fgHiRed.Printf("It's over than sector Size(%d) entered: %d\n", maxSectorSize, idx)
				fgHiRed.Printf("Reseted maximum size: %d\n", maxSectorSize)
				printBlock(content, maxSectorSize*sectorSize, contentSize, []foundHexArray{})
			} else if idx == maxSectorSize {
				lastStartIdx = idx * sectorSize
				lastEndIdx = idx * sectorSize

				clear()
				printBlock(content, lastStartIdx, lastEndIdx, []foundHexArray{})
			} else {
				lastStartIdx = idx * sectorSize
				lastEndIdx = (idx + 1) * sectorSize

				clear()
				printBlock(content, lastStartIdx, lastEndIdx, []foundHexArray{})
			}
		}
	}

}

func findHexArray(arr []byte, start int64, end int64, matchStr []int) []foundHexArray {
	hexLen := int64(len(matchStr))

	var hexArr []byte
	for i := int64(0); i < hexLen; i++ {
		hexArr = append(hexArr, byte(matchStr[i]))
	}

	var resultSets []foundHexArray
	for i := start; i < end; i++ {
		if bytes.Compare(arr[i:i+hexLen], hexArr) == 0 {
			fmt.Printf("matched start: %d, end: %d\n", i, i+hexLen)
			resultSets = append(resultSets, foundHexArray{
				startIdx: i,
				endIdx:   i + hexLen,
			})
		}
	}
	if len(resultSets) == 0 {
		return []foundHexArray{
			{
				startIdx: -1,
				endIdx:   -1,
			},
		}
	}
	return resultSets

}

func showSectorInformation() {
	fileName := readFilePath()
	var err error
	content, err = os.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Cannot open file: " + fileName)
	}
	contentSize = int64(len(content))
	findSector()
}

func printBlock(arr []byte, start int64, end int64, highlightHexArraySets []foundHexArray) {

	hlHexArrIdx := 0
	hlHexArrSize := len(highlightHexArraySets)

	var rowWidth int64 = 16
	var currRowWidth int64 = 16
	var seenBytes int64 = start

	fgHiGreenUnder.Printf("Total byte size: %d\n", contentSize)
	fgHiGreenUnder.Printf("Total sector size: %d\n", contentSize/sectorSize)
	fmt.Printf("Current Byte Idx: %d\n", start)
	fgHiCyan.Printf(" offset(Dec|	   Hex)  ")
	fmt.Printf(" 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n")
	fmt.Printf(" ===================================================================================================\n")

	for i := start; i < end; i += rowWidth {
		if (int64(len(arr)) - seenBytes) < rowWidth {
			currRowWidth = int64(len(arr)) - seenBytes
		}
		row := arr[i:(seenBytes + currRowWidth)]

		fmt.Printf("%10d| %10X  | ", seenBytes, seenBytes)

		for i := int64(0); i < rowWidth; i++ {
			if i < currRowWidth {
				if hlHexArrSize > hlHexArrIdx && (i >= highlightHexArraySets[hlHexArrIdx].startIdx-seenBytes &&
					i < highlightHexArraySets[hlHexArrIdx].endIdx-seenBytes) {
					fgHiCyanUnder.Printf("%02x", row[i])
					fmt.Printf(" ")
					if (i + 1) == highlightHexArraySets[hlHexArrIdx].endIdx-seenBytes {
						hlHexArrIdx++
					}
				} else {
					fmt.Printf("%02x ", row[i])
				}
			} else {
				fmt.Print(strings.Repeat(" ", 3))
			}
		}

		fmt.Print("|")
		fmt.Print(" ")

		for i := int64(0); i < rowWidth; i++ {
			if i < currRowWidth {
				if row[i] >= 0x20 && row[i] < 0x7f {
					fmt.Print(string(row[i]))
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Print(strings.Repeat(" ", 3))
			}
		}
		fmt.Print("|")
		fmt.Print("\n")

		seenBytes += rowWidth
	}
	fmt.Printf("Current sector idx: ")
	fgHiCyanUnder.Printf("%d", start/int64(sectorSize))
	fmt.Printf("/%d\n\n\n", contentSize/sectorSize)
}

func showPartitionInformation() {

}
