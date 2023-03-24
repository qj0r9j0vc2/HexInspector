package headerParse

import (
	"HexInspector/util"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func (elfFs *ELFFile) setArch() {
	switch elf.Class(elfFs.Ident[elf.EI_CLASS]) {
	case elf.ELFCLASS64:
		elfFs.Header = new(elf.Header64)
		elfFs.FileHdr.Arch = elf.ELFCLASS64

	case elf.ELFCLASS32:
		elfFs.Header = new(elf.Header32)
		elfFs.FileHdr.Arch = elf.ELFCLASS32
	default:
		fmt.Println("Elf Arch Class Invalid !")
		os.Exit(1)
	}
}

func (elfFs *ELFFile) mapHeader() {

	switch elf.Data(elfFs.Ident[elf.EI_DATA]) {
	case elf.ELFDATA2LSB:
		elfFs.FileHdr.Endianness = binary.LittleEndian
	case elf.ELFDATA2MSB:
		elfFs.FileHdr.Endianness = binary.BigEndian
	default:
		fmt.Println("Possible Corruption, Endianness unknown")
	}

	elfFs.Fh.Seek(0, io.SeekStart)
	err := binary.Read(elfFs.Fh, elfFs.FileHdr.Endianness, elfFs.Header)
	checkError(err)

	switch h := elfFs.Header.(type) {
	case *elf.Header32:
		elfFs.FileHdr.Machine = elf.Machine(h.Machine)
	case *elf.Header64:
		elfFs.FileHdr.Machine = elf.Machine(h.Machine)
	}
}

func printELFHeader(hdr interface{}) {
	if h, ok := hdr.(*elf.Header64); ok {
		fmt.Printf("-------------------------- Elf Header ------------------------\n")
		fmt.Printf("Magic: % x\n", h.Ident)
		fmt.Printf("Class: %s\n", elf.Class(h.Ident[elf.EI_CLASS]))
		fmt.Printf("Data: %s\n", elf.Data(h.Ident[elf.EI_DATA]))
		fmt.Printf("Version: %s\n", elf.Version(h.Version))
		fmt.Printf("OS/ABI: %s\n", elf.OSABI(h.Ident[elf.EI_OSABI]))
		fmt.Printf("ABI Version: %d\n", h.Ident[elf.EI_ABIVERSION])
		fmt.Printf("Elf Type: %s\n", elf.Type(h.Type))
		fmt.Printf("Machine: %s\n", elf.Machine(h.Machine))
		fmt.Printf("Entry: 0x%x\n", h.Entry)
		fmt.Printf("Program Header Offset: 0x%x\n", h.Phoff)
		fmt.Printf("Section Header Offset: 0x%x\n", h.Shoff)
		fmt.Printf("Flags: 0x%x\n", h.Flags)
		fmt.Printf("Elf Header Size (bytes): %d\n", h.Ehsize)
		fmt.Printf("Program Header Entry Size (bytes): %d\n", h.Phentsize)
		fmt.Printf("Number of Program Header Entries: %d\n", h.Phnum)
		fmt.Printf("Size of Section Header Entry: %d\n", h.Shentsize)
		fmt.Printf("Number of Section Header Entries: %d\n", h.Shnum)
		fmt.Printf("Index of section header string table: %d\n", h.Shstrndx)
	}

	if h, ok := hdr.(*elf.Header32); ok {
		fmt.Printf("-------------------------- Elf Header ------------------------\n")
		fmt.Printf("Magic: % x\n", h.Ident)
		fmt.Printf("Class: %s\n", elf.Class(h.Ident[elf.EI_CLASS]))
		fmt.Printf("Data: %s\n", elf.Data(h.Ident[elf.EI_DATA]))
		fmt.Printf("Version: %s\n", elf.Version(h.Version))
		fmt.Printf("OS/ABI: %s\n", elf.OSABI(h.Ident[elf.EI_OSABI]))
		fmt.Printf("ABI Version: %d\n", h.Ident[elf.EI_ABIVERSION])
		fmt.Printf("Elf Type: %s\n", elf.Type(h.Type))
		fmt.Printf("Machine: %s\n", elf.Machine(h.Machine))
		fmt.Printf("Entry: 0x%x\n", h.Entry)
		fmt.Printf("Program Header Offset: 0x%x\n", h.Phoff)
		fmt.Printf("Section Header Offset: 0x%x\n", h.Shoff)
		fmt.Printf("Flags: 0x%x\n", h.Flags)
		fmt.Printf("Elf Header Size (bytes): %d\n", h.Ehsize)
		fmt.Printf("Program Header Entry Size (bytes): %d\n", h.Phentsize)
		fmt.Printf("Number of Program Header Entries: %d\n", h.Phnum)
		fmt.Printf("Size of Section Header Entry: %d\n", h.Shentsize)
		fmt.Printf("Number of Section Header Entries: %d\n", h.Shnum)
		fmt.Printf("Index of section header string table: %d\n", h.Shstrndx)
	}
	return
}

func showELFHeader() {
	fileName := util.ReadFilePath()

	var file ELFFile
	file.Fh, file.Err = os.Open(fileName)
	if file.Err != nil {
		fmt.Println("Cannot Open file: " + file.Err.Error())
		return
	}
	file.Fh.Read(file.Ident[:16])

	if isElf(file.Ident[:4]) == false {
		fmt.Println("This is not an Elf binary")
		os.Exit(1)
	}
	file.setArch()
	file.mapHeader()

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	util.Clear()
	fmt.Println(file.Fh)
	go printELFHeader(file.Header)
	<-stopper
}

func isElf(magic []byte) bool {
	return !(magic[0] != '\x7f' || magic[1] != 'E' || magic[2] != 'L' || magic[3] != 'F')
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
