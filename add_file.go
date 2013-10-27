package main

import (
  "bufio"
  "github.com/mreiferson/go-snappystream"
  "crypto/sha256"
  "encoding/gob"
  "io"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
) 

func cmd_add_files(archive_filename string, infiles []string) (error) {

  // Create a temp directory for processing the files
  tempdir, err := ioutil.TempDir("", "fastpack")
  log.Printf("Using tempdir: %s", tempdir)
  if err != nil { return err }
  defer os.RemoveAll(tempdir)

  //Open the output file
  archive, err := os.Create(archive_filename)
  if err != nil {
    return err
  }
  defer archive.Close()

  // A place for all of the headers
  var headers []Filenode

  // Time to rip through the files
  for index, element := range infiles {
    tempfile, file_header := process_infile(index, element)
    headers = append(headers, file_header)

    log.Printf("Packing %s", element)
  }

  // Let's encode those headers into a gob so we can store them
  enc := gob.NewEncoder(archive)
  err = enc.Encode(headers)
  if err != nil {
    return err
  }

  return nil
}

func process_infile(index int, filename string) (*os.File, Filenode) {
  log.Printf("Processing %s", filename)
  // Build up some information about this file
  var file_header Filenode
  file_header.Name = filepath.Base(filename)
  file_header.Path = filepath.Dir(filename)

  file_stat, err := os.Stat(filename)
  if err != nil { log.Fatal(err) }
  file_header.Size = file_stat.Size()
  file_header.mode = file_stat.Mode()
  file_header.modification_time = file_stat.ModTime()
  file_header.Isdir = file_stat.IsDir()

  // Open our input file
  infile, err := os.Open(filename)
  if err != nil { log.Fatal(err) }
  defer infile.Close()

  // Create the tempfile
  tempfile, err := ioutil.TempFile(tempdir, file_header.Name)
  if err != nil { log.Fatal(err) }
  defer tempfile.Close()
  
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
    if err != nil && err != io.EOF { log.Fatal(err) }
    // If the returned size is zero, we're at the end of the file
    if read_count == 0 { break }

    // Add the buffer contents to the checksum calculation
    checksum.Write(buf[:read_count])

    // And write the buffer to the output stream
    if _, err := snappyWriter.Write(buf[:read_count]); err != nil {
      log.Fatal(err)
    }
  }

  // Flush anything left in the final chain of the output writer
  if err = writer.Flush(); err != nil { log.Fatal(err) }

  file_header.Checksum = checksum.Sum(nil)

  return tempfile, file_header
}
