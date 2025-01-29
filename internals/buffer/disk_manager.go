package main

import (
	"errors"
)

const DiskMaxNumPages = 15

type DiskManager interface {
	ReadPage(PageID) (*Page, error)
	WritePage(*Page) error
	AllocatePage() *PageID
	DeallocatePage(PageID)
}

type DiskManagerMock struct {
	numPage int 
	pages   map[PageID]*Page
}

func (d *DiskManagerMock) ReadPage(pageID PageID) (*Page, error) {
	if page, ok := d.pages[pageID]; ok {
		return page, nil
	}

	return nil, errors.New("Page not found")
}

func (d *DiskManagerMock) WritePage(page *Page) error {
	d.pages[page.id] = page
	return nil
}

func (d *DiskManagerMock) AllocatePage() *PageID {
	if d.numPage == DiskMaxNumPages-1 {
		return nil
	}
	d.numPage = d.numPage + 1
	pageID := PageID(d.numPage)
	return &pageID
}

func (d *DiskManagerMock) DeallocatePage(pageID PageID) {
	delete(d.pages, pageID)
}

func NewDiskManagerMock() *DiskManagerMock {
	return &DiskManagerMock{-1, make(map[PageID]*Page)}
}