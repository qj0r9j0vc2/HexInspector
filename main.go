package main

import (
	"HexInspector/SectorInformation"
	"HexInspector/fileOpen"
	"HexInspector/util"
	"github.com/manifoldco/promptui"
	"log"
)

func main() {
	menu()
}

func menu() {
	util.Clear()
	for {
		menu := []string{"File Open", "Show Sector Information", "Show Partition Information", "Exit"}
		prompt := promptui.Select{Label: "Select Menu", Items: menu}
		i, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		switch i {
		case 0:
			fileOpen.FileOpen()
			break
		case 1:
			SectorInformation.ShowSectorInformation()
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

func showPartitionInformation() {

}
