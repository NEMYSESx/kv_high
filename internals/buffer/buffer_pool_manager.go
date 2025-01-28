package main

const MaxPoolSize = 10

type BufferPoolManager struct {
	diskManager DiskManager
	pages       [MaxPoolSize]*Page
	replacer    *ClockReplacer
	freeList    []FrameID
	pageTable   map[PageID]FrameID
}
