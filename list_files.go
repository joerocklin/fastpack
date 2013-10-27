package main

import (
  "bufio"
  "encoding/gob"
  "fmt"
  "log"
  "os"
)

func cmd_list_files(archive_filename string) (error) {
  // Open our input file
  archive, err := os.Open(archive_filename)
  if err != nil { log.Fatal(err) }
  defer archive.Close()

  // Create the read buffer
  reader := bufio.NewReader(archive)

  // Create the gob decoder
  dec := gob.NewDecoder(reader)

  var headers []Filenode
  err = dec.Decode(&headers)
  if err != nil { return err }

  for _, file := range headers {
    fmt.Printf("%d     %s/%s\n", file.Size, file.Path, file.Name)
  }

  return nil
}