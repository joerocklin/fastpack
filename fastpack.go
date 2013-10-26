package main

import (
  "bufio"
  "github.com/mreiferson/go-snappystream"
  "crypto/sha256"
  "flag"
  "fmt"
  "io"
  "log"
  "os"
  "time"
)

//type encryption_type 

type filenode struct {
  name string
  path string
  size int64
 
  mode os.FileMode
  isdir bool
  modification_time time.Time
  uid uint16
  gid uint16

  checksum [sha256.Size]byte
}

var help bool
var outfile string
var infiles []string
var logger = log.New(os.Stderr, "fastpack", log.Flags())

func init() {
  flag.BoolVar(&help, "h", false, "help")
  flag.Parse()

  args := flag.Args()
  
  if len(args) > 1 { 
    outfile = args[0]
    infiles = args[1:]
  } else {
    help = true
  }
}

func process_infile(index int, filename string) {
  infile, err := os.Open(filename)
  if err != nil { log.Fatal(err) }
  defer infile.Close()

  outfile, err := os.Create(fmt.Sprintf("%s.snap", filename))
  if err != nil { log.Fatal(err) }
  defer outfile.Close()

  reader := bufio.NewReader(infile)
  writer := bufio.NewWriter(outfile)
  snappyWriter := snappystream.NewWriter(writer)

  checksum := sha256.New()

  buf := make([]byte, 1024)
  for {
    read_count, err := reader.Read(buf)
    if err != nil && err != io.EOF { log.Fatal(err) }
    if read_count == 0 { break }

    checksum.Write(buf[:read_count])
    if _, err := snappyWriter.Write(buf[:read_count]); err != nil {
      log.Fatal(err)
    }

  }

  if err = writer.Flush(); err != nil { log.Fatal(err) }

  log.Println(checksum.Sum(nil))
}

func main() {
  if help {
    flag.Usage()
    return;
  }

  log.Println("Writing output to", outfile)

  for index, element := range infiles {
    process_infile(index, element)
  }
}