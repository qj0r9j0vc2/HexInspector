package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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

func main() {
	fileName := os.Args[1]
	if fileName == "" {
		log.Fatalln("fileName not input")
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Cannot open file: " + fileName)
	}

	printBlock(content)
}
