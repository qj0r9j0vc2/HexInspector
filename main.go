package main

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strings"
)

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

func showSectorInformation() {
	fileName := readFilePath()
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Cannot open file: " + fileName)
	}

	printBlock(content)
}

func printBlock(arr []byte) {
	rowWidth := 16
	currRowWidth := 16
	seenBytes := 0

	for i := 0; i < len(arr); i += rowWidth {
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
