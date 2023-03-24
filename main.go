package main

import (
	"HexInspector/fileOpen"
	"HexInspector/headerParse"
	"HexInspector/showHex"
	"HexInspector/util"
	"github.com/manifoldco/promptui"
	"log"
)

func main() {
	menu()
}

func menu() {
	util.Clear()
	menu := []string{"File Open", "Show Sector Information", "Show File Header", "Exit"}
	prompt := promptui.Select{Label: "Select Menu", Items: menu}
	for {
		i, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		switch i {
		case 0:
			fileOpen.FileOpen()
			break
		case 1:
			showHex.ShowHex()
			break
		case 2:
			headerParse.PrintFileSystem()
			break
		case 3:
			println("Exit..!!")
			return
		}
	}
}
