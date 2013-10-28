package main

import (
	"encoding/binary"
  "os"
  "time"
)

const fastpackIndexMagicv1 = 0x667061636b303031

const indexOffsetSize = binary.MaxVarintLen64

type Fileindex struct {
	Name        string
	Path        string
	Offset      int64
	Size        int64
	Expand_Size int64

	Mode    os.FileMode
	Isdir   bool
	ModTime time.Time
	Uid     uint16
	Gid     uint16

	Checksum []byte
}

func IntToBytes(v int64) []byte {
	buf := make([]byte, indexOffsetSize)

	binary.PutVarint(buf, v)

	return buf
}

func BytesToInt(buf []byte) int64 {
	val, _ := binary.Varint(buf)

	return val
}
