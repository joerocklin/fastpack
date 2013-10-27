package main

import (
	"os"
	"time"
)

type Filenode struct {
	Name string
	Path string
	Size int64

	Mode    os.FileMode
	Isdir   bool
	ModTime time.Time
	Uid     uint16
	Gid     uint16

	Checksum []byte
}
