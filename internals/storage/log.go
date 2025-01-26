package main

import (
	"fmt"
	"hash/crc32"
)

type Record struct {
    CRC       uint32
    Timestamp uint32  
    KeySize   uint32
    ValueSize uint32
    Key       string
    Value     string
}

func main(){
	data := []byte("Hello world!")

	fmt.Println(data)

	checkSum := crc32.ChecksumIEEE(data)

	fmt.Printf("CRC-32 (IEEE): %08x\n", checkSum)
}