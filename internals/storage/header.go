package storage

//HeaderBuilder groups all logic related to header creation, parsing, and validation into one place
// hb := NewHeaderBuilder()
// hb.BuildHeader(...)
// hb.ParseHeader(...)
// hb.ValidateHeader(...)

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

type Header struct {
	CRC       uint32
	Length    uint16
	BlockType byte
}

type HeaderBuilder struct{}

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{}
}

func (hb *HeaderBuilder) BuildHeader(data []byte, blockType int) ([]byte, error) {
	crc := crc32.ChecksumIEEE(data);
	//packing - turning data into raw bytes
	//we doing this couz Storing binary data or sending data over the network
	headerBuf := make([]byte,recordSize)
	binary.LittleEndian.PutUint32(headerBuf[0:4],crc)
	binary.LittleEndian.PutUint16(headerBuf[4:6],uint16(len(data)))
	headerBuf[6] = byte(blockType)
	
	return headerBuf,nil
}

func (hb *HeaderBuilder) ParseHeader(headerBuf []byte) (Header,error){
	if len(headerBuf) != recordHeaderSize {
		return Header{}, fmt.Errorf("The header size doesnot matches: expected %d, got %d",recordHeaderSize,len(headerBuf))
	}
	//un packing
	header := Header{
		CRC: binary.LittleEndian.Uint32(headerBuf[0:4]),
		Length: binary.LittleEndian.Uint16(headerBuf[4:6]),
		BlockType: byte(headerBuf[6]),
	}

	return header,nil
}

func (hb *HeaderBuilder) ValidateHeader(header Header,data []byte) error {
	expectedCRC := crc32.ChecksumIEEE(data)
	if header.CRC != expectedCRC {
		return fmt.Errorf("CRC mismatche: expected %d, got %d", expectedCRC,header.CRC)
	}
	return nil
}