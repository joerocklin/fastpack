package main

import (
	//"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func cmd_list_files(archive_filename string) error {
	// Open our input file
	archive, err := os.Open(archive_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	// Create the read buffer
	//reader := bufio.NewReader(archive)

	headers, err := getFileIndex(archive)

  fmt.Printf("%-8s %-16s %s\n", "offset", "Size (bytes)", "Path/File")
	for _, file := range headers {
		fmt.Printf("%-8d %-16d %s/%s\n", file.Offset, file.Size, file.Path, file.Name)
	}

	return nil
}

func getFileIndex(archive *os.File) ([]Fileindex, error) {
	// Store the current file position
	curpos, _ := archive.Seek(0, os.SEEK_CUR)

	// Seek to the end of the file, and read in the value which tells us where
	// the start of the index is located
	archive.Seek(-indexOffsetSize, os.SEEK_END)

	offset_bytes := make([]byte, indexOffsetSize)
	_, err := archive.Read(offset_bytes)
	offset := BytesToInt(offset_bytes)

	log.Printf("Index offset located at %d bytes from end", offset)

	archive.Seek(-(indexOffsetSize + offset), os.SEEK_END)

	// Create the gob decoder
	dec := gob.NewDecoder(archive)

	var headers []Fileindex
	err = dec.Decode(&headers)
	if err != nil {
		return nil, err
	}

	// restore the original file position
	archive.Seek(curpos, os.SEEK_SET)

	return headers, nil
}
