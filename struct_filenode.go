package main

import (
  "os"
  "time"
) 

type Filenode struct {
  Name string
  Path string
  Size int64
 
  mode os.FileMode
  Isdir bool
  modification_time time.Time
  Uid uint16
  Gid uint16

  Checksum []byte
}