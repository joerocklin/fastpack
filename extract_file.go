package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	//	"github.com/mreiferson/go-snappystream"
	//  "errors"
	//  "fmt"
	"io"
	"log"
	"os"
)

func cmd_extract_files(archive_filename string, files []string) error {
	// Open our input file
	archive, err := os.Open(archive_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	// Create the read buffer
	//reader := bufio.NewReader(archive)

	headers, err := getFileIndex(archive)
	if err != nil {
		return err
	}

	for _, element := range headers {
		err := process_outfile(archive_filename, element)
		if err != nil {
			return err
		}
	}

	return nil
}

func process_outfile(archive_filename string, index Fileindex) error {
	log.Printf("Processing %s", index.Name)

	// Open our input file
	archive, err := os.Open(archive_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()
	file_reader := io.LimitReader(archive, 100)

	// Build our decoder
	dec := gob.NewDecoder(file_reader)
	var file_header Fileindex

	// Seek to the start of this file
	log.Printf("Seeking to %x", index.Offset)
	_, err = archive.Seek(index.Offset, os.SEEK_SET)
	if err != nil {
		return err
	}

	// Read in the header
	err = dec.Decode(&file_header)
	if err != nil {
		return err
	}

	log.Printf("%+v", file_header)
	pos, _ := archive.Seek(0, os.SEEK_CUR)
	log.Printf("Current pos: %x", pos)

	log.Printf("Limiting reads to %d", index.Size)
	//file_reader := io.LimitReader(archive, index.Size)
	file_reader = archive
	//snappyReader := snappystream.NewReader(file_reader, snappystream.VerifyChecksum)
	reader := file_reader

	// Create the outfile
	outfile, err := os.Create(file_header.Name)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(outfile)

	// Start a checksum for the original file
	checksum := sha256.New()

	// Now do the work: read from the input chain, and write to the output file
	buf := make([]byte, 4096)
	for {
		read_count, err := reader.Read(buf)
		// If we get an error that is not EOF, then we have a problem
		if err != nil && err != io.EOF {
			log.Println("Error reading")
			return err
		}
		log.Printf("Read %d bytes from stream", read_count)
		// If the returned size is zero, we're at the end of the file
		if read_count == 0 {
			log.Println(err)
			break
		}

		// Add the buffer contents to the checksum calculation
		checksum.Write(buf[:read_count])

		// And write the buffer to the output stream
		if _, err := writer.Write(buf[:read_count]); err != nil {
			return err
		}
	}

	// Flush anything left in the final chain of the output writer
	if err = writer.Flush(); err != nil {
		log.Fatal(err)
	}

	if !bytes.Equal(file_header.Checksum, checksum.Sum(nil)) {
		log.Printf("     Got: %+v", checksum.Sum(nil))
		log.Printf("Expected: %+v", file_header.Checksum)
		//return errors.New(fmt.Sprintf("Checksum mismatch for file %s/%s", file_header.Path, file_header.Name))
	}

	return nil
}
