package headerParse

import (
	"HexInspector/util"
	"github.com/manifoldco/promptui"
	"log"
)

func PrintFileSystem() {
	util.Clear()
	menu := []string{"ELF", "FAT32", "Back"}
	prompt := promptui.Select{Label: "Select FileSystem", Items: menu}
	for {
		i, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		switch i {
		case 0:
			showELFHeader()
			break
		case 1:
			showFATHeader()
			break
		case 2:
			return
		}
	}

}
