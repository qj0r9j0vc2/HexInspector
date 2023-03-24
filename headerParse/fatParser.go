package headerParse

import (
	"HexInspector/util"
	"fmt"
	"github.com/mitchellh/go-fs"
	"github.com/mitchellh/go-fs/fat"
	"os"
	"os/signal"
	"syscall"
)

func printFATHeader(sector fat.BootSectorCommon) {
	fmt.Printf("-------------------------- FAT Header ------------------------\n")
	fmt.Printf("Byte Per Sector: %08x(%d)\n", sector.BytesPerSector, sector.BytesPerSector)
	fmt.Printf("Byte Per Cluster: %d\n", uint16(sector.SectorsPerCluster)*sector.BytesPerSector)
	fmt.Printf("Reserved Sector Count: %d\n", sector.ReservedSectorCount)
	fmt.Printf("Total Sector: %d\n", sector.TotalSectors)
	fmt.Printf("FAT Size: %d\n", sector.SectorsPerFat)
	return
}

func showFATHeader() {
	fileName := util.ReadFilePath()

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Cannot Open file: " + err.Error())
		return
	}
	disk, err := fs.NewFileDisk(f)
	if err != nil {
		return
	}
	sector, err := fat.DecodeBootSector(disk)
	if err != nil {
		return
	}

	//for i := 1; i < int(sector.NumFATs); i++ {
	//	decoded, err := fat.DecodeFAT(disk, sector, i)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//
	//	fat.DecodeBootSector()
	//}

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	util.Clear()
	printFATHeader(*sector)
	<-stopper
}
