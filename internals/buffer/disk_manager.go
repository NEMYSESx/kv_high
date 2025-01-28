package main

const DiskMaxNumPages = 15

type DiskManager interface {
	ReadPage(PageID) (*Page, error)
	WritePage(*Page) error
	AllocatePage() *PageID
	DeallocatePage(PageID)
}