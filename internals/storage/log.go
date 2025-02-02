package main

import (
	"fmt"
	"hash/crc32"
	"os"
)

type Record struct {
	CRC       uint32
	Timestamp uint32  
	KeySize   uint32
	ValueSize uint32
	Key       string
	Value     string
}

func main() {
	file, err := os.OpenFile("append.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("New log entry\n")
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}
	fmt.Println("File written successfully!")

	data := []byte("Hello world!")
	fmt.Println("Data:", string(data))

	checkSum := crc32.ChecksumIEEE(data)
	fmt.Printf("CRC-32 (IEEE): %08x\n", checkSum)
}
