package main

import (
	"bufio"
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
var sectorSize = 512
var fgHiGreen = color.New(color.FgHiGreen).Add(color.Underline)

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

func inputInt() int {
	str := inputStr()
	integer, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("입력된 데이터가 올바르지 않습니다.: %s \n", err.Error())
	}
	return integer
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

func findSector() {
	printBlock(content, 0, sectorSize)
	for {
		fgHiGreen.Println("Enter sector idx: ")
		idx := inputInt()

		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()

		printBlock(content, idx*sectorSize, idx*sectorSize+sectorSize)

	}
	
}

func showSectorInformation() {
	fileName := readFilePath()
	var err error
	content, err = os.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Cannot open file: " + fileName)
	}

	findSector()
}

func printBlock(arr []byte, start int, end int) {
	rowWidth := 16
	currRowWidth := 16
	seenBytes := start

	fgHiGreen.Printf("Total byte size: %d\n", len(arr))
	fgHiGreen.Printf("Total sector size: %d\n", len(arr)/sectorSize)
	fmt.Printf("Current Byte Idx: %d\n", start)
	fmt.Printf(" offset(h)			  ")
	fmt.Printf(" 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F\n")
	fmt.Printf(" ===================================================================================================\n")

	for i := start; i < end; i += rowWidth {
		if (len(arr) - seenBytes) < rowWidth {
			currRowWidth = len(arr) - seenBytes
		}
		row := arr[i:(seenBytes + currRowWidth)]

		fmt.Printf("Dec/Hex: %10d| %10X  | ", seenBytes, seenBytes)

		for i := 0; i < rowWidth; i++ {
			if i < currRowWidth {
				fmt.Printf("%02x ", row[i])
			} else {
				fmt.Print(strings.Repeat(" ", 3))
			}
		}

		fmt.Print("|")
		fmt.Print(" ")

		for i := 0; i < rowWidth; i++ {
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
}

func showPartitionInformation() {

}
