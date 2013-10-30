package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"github.com/mreiferson/go-snappystream"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func cmd_add_files(archive_filename string, infiles []string) error {
	// Create a temp directory for processing the files
	tempdir, err := ioutil.TempDir("", "fastpack")
	log.Printf("Using tempdir: %s", tempdir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempdir)

	//Open the output file
	archive, err := os.Create(archive_filename)
	if err != nil {
		return err
	}
	defer archive.Close()

	// A place for all of the file info
	var archive_index []Fileindex

	// Time to rip through the files
	for _, element := range infiles {
		tempfile, file_header := process_infile(element)
		defer tempfile.Close()

		tempfile.Seek(0, os.SEEK_SET)

		var file_index Fileindex
		file_index.Name = file_header.Name
		file_index.Path = file_header.Path
		file_stat, err := tempfile.Stat()
		if err != nil {
			return err
		}
		file_index.Size = file_stat.Size()

		file_index.Offset, err = archive.Seek(0, os.SEEK_CUR)
		if err != nil {
			return err
		}

		log.Printf("Adding %s", element)

		// Create the read buffer
		reader := io.Reader(tempfile)

		// Chain the output buffers
		writer := io.Writer(archive)

		buf := make([]byte, 1024)
		for {
			read_count, err := reader.Read(buf)
			// If we get an error that is not EOF, then we have a problem
			if err != nil && err != io.EOF {
				return err
			}

			// If the returned size is zero, we're at the end of the file
			if read_count == 0 {
				break
			}

			// And write the buffer to the output stream
			_, err = writer.Write(buf[:read_count])
			if err != nil {
				return err
			}
		}

		// // Flush anything left in the final chain of the output writer
		// if err = writer.Flush(); err != nil {
		// 	log.Fatal(err)
		// }
		log.Printf("%+v", file_index)
		archive_index = append(archive_index, file_index)
	}

	archive.Write(IntToBytes(fastpackIndexMagicv1))
	pre_index_position, err := archive.Seek(0, os.SEEK_CUR)

	// Let's encode those headers into a gob so we can store them
	enc := gob.NewEncoder(archive)
	err = enc.Encode(archive_index)
	if err != nil {
		return err
	}
	post_index_position, err := archive.Seek(0, os.SEEK_CUR)

	var Offset int64
	Offset = post_index_position - pre_index_position
	log.Printf("Index offset located at %d bytes from end", Offset)
	archive.Write(IntToBytes(Offset))

	return nil
}

func process_infile(filename string) (*os.File, Fileindex) {
	log.Printf("Processing %s", filename)
	// Build up some information about this file
	var file_header Fileindex
	file_header.Name = filepath.Base(filename)
	file_header.Path = filepath.Dir(filename)

	file_stat, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	file_header.Expand_Size = file_stat.Size()
	file_header.Mode = file_stat.Mode()
	file_header.ModTime = file_stat.ModTime()
	file_header.Isdir = file_stat.IsDir()

	// Open our input file
	infile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// Create the tempfile
	tempfile, err := ioutil.TempFile(tempdir, file_header.Name)
	if err != nil {
		log.Fatal(err)
	}

	// Create the read buffer
	reader := bufio.NewReader(infile)

	// Chain the output buffers
	writer := bufio.NewWriter(tempfile)
	snappyWriter := snappystream.NewWriter(writer)

	// Start a checksum for the original file
	checksum := sha256.New()

	// Now do the work: read from the input buffer, and write to the output chain
	buf := make([]byte, 1024)
	for {
		read_count, err := reader.Read(buf)
		// If we get an error that is not EOF, then we have a problem
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		// If the returned size is zero, we're at the end of the file
		if read_count == 0 {
			break
		}

		// Add the buffer contents to the checksum calculation
		checksum.Write(buf[:read_count])

		// And write the buffer to the output stream
		if _, err := snappyWriter.Write(buf[:read_count]); err != nil {
			log.Fatal(err)
		}
	}

	// Flush anything left in the final chain of the output writer
	if err = writer.Flush(); err != nil {
		log.Fatal(err)
	}

	file_header.Checksum = checksum.Sum(nil)
	log.Printf("%+v", file_header)

	return tempfile, file_header
}
