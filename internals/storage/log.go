package storage

type RecordType struct {
	ZeroType   int
	FullType   int
	FirstType  int
	MiddleType int
	LastType   int
}

var record = &RecordType{
	ZeroType:   0,
	FullType:   1,
	FirstType:  2,
	MiddleType: 3,
	LastType:   4,
}

const recordSize = 32768

const recordHeaderSize = 4 + 2 + 1
