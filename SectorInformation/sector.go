package SectorInformation

import (
	"HexInspector/util"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FoundHexArray struct {
	startIdx int64
	endIdx   int64
}

var (
	content     []byte
	sectorSize  = int64(512)
	contentSize int64
)

func ShowSectorInformation() {
	fileName := util.ReadFilePath()
	var err error
	content, err = os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Cannot open file: " + fileName)
		fmt.Println(err.Error())
	} else {
		contentSize = int64(len(content))
		findSector()
	}
}

func findSector() {
	lastStartIdx := int64(0)
	lastEndIdx := sectorSize
	maxsectorSize := contentSize / sectorSize

	util.Clear()
	printBlock(content, 0, int64(sectorSize), []FoundHexArray{})

	for {
		util.FgHiGreenUnder.Println("Enter sector idx")
		fmt.Printf("(if you want to search hex value, enter hex value after 'f')\n")
		str := util.InputStr()
		idx, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			if str == "f" || str == "F" {
				var str string
				str = util.InputStr()
				str = strings.ReplaceAll(str, " ", "")

				FoundHexArray := findHexArray(content, lastStartIdx, lastEndIdx, util.StringToHex(str))

				util.Clear()
				printBlock(content, lastStartIdx, lastEndIdx, FoundHexArray)
			}
		} else {
			if idx > maxsectorSize {
				util.Clear()

				util.FgHiRed.Printf("It's over than sector Size(%d) entered: %d\n", maxsectorSize, idx)
				util.FgHiRed.Printf("Reseted maximum size: %d\n", maxsectorSize)
				printBlock(content, maxsectorSize*sectorSize, contentSize, []FoundHexArray{})
			} else if idx == maxsectorSize {
				lastStartIdx = idx * sectorSize
				lastEndIdx = idx * sectorSize

				util.Clear()
				printBlock(content, lastStartIdx, lastEndIdx, []FoundHexArray{})
			} else {
				lastStartIdx = idx * sectorSize
				lastEndIdx = (idx + 1) * sectorSize

				util.Clear()
				printBlock(content, lastStartIdx, lastEndIdx, []FoundHexArray{})
			}
		}
	}

}

func findHexArray(arr []byte, start int64, end int64, matchStr []int) []FoundHexArray {
	hexLen := int64(len(matchStr))

	var hexArr []byte
	for i := int64(0); i < hexLen; i++ {
		hexArr = append(hexArr, byte(matchStr[i]))
	}

	var resultSets []FoundHexArray
	for i := start; i < end; i++ {
		if bytes.Compare(arr[i:i+hexLen], hexArr) == 0 {
			fmt.Printf("matched start: %d, end: %d\n", i, i+hexLen)
			resultSets = append(resultSets, FoundHexArray{
				startIdx: i,
				endIdx:   i + hexLen,
			})
		}
	}
	if len(resultSets) == 0 {
		return []FoundHexArray{
			{
				startIdx: -1,
				endIdx:   -1,
			},
		}
	}
	return resultSets

}

func printBlock(arr []byte, start int64, end int64, highlightHexArraySets []FoundHexArray) {

	hlHexArrIdx := 0
	hlHexArrSize := len(highlightHexArraySets)

	var rowWidth int64 = 16
	var currRowWidth int64 = 16
	var seenBytes int64 = start

	util.FgHiGreenUnder.Printf("Total byte size: %d\n", contentSize)
	util.FgHiGreenUnder.Printf("Total sector size: %d\n", contentSize/sectorSize)
	fmt.Printf("Current Byte Idx: %d\n", start)
	util.FgHiCyan.Printf(" offset(Dec|	   Hex)  ")
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
					util.FgHiCyanUnder.Printf("%02x", row[i])
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
	util.FgHiCyanUnder.Printf("%d", start/int64(sectorSize))
	fmt.Printf("/%d\n\n\n", contentSize/sectorSize)
}
