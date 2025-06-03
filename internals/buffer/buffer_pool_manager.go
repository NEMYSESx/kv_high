package main

import (
	"errors"
)

const MaxPoolSize = 4

// pages is an array of pointers to Page.
// Each slot in the array corresponds to a frame in memory.
// To access a page, we need a frameID (which is just an index in this array).

// This pageTable is a map (hash table) that stores PageID → FrameID mappings.
// Key: PageID (the ID of a page requested by the user).
// Value: FrameID (the slot in the pages array where that page is stored).

// What is frameID?
// frameID is just an index in the pages array.
// It represents a physical slot in the buffer pool where a Page is stored.
// It is not tied to any specific page permanently. A frameID can hold different pages over time.

// Why do we need frameID instead of directly mapping pageID to pages?
// pages is a fixed-size array (e.g., MaxPoolSize = 4).
// The database can have thousands or millions of pages (PageIDs).
// We cannot create an array with millions of slots just for PageIDs because it would be inefficient.
// 👉 Instead, we use frameID to manage a limited number of slots dynamically.

type BufferPoolManager struct {
	diskManager DiskManager
	pages       [MaxPoolSize]*Page
	replacer    *ClockReplacer
	freeList    []FrameID
	pageTable   map[PageID]FrameID
}

func (b *BufferPoolManager) FetchPage(pageID PageID) *Page {
	if frameID, ok := b.pageTable[pageID]; ok {
		page := b.pages[frameID]
		page.pinCount++
		(*b.replacer).Pin(frameID)
		return page
	}

	//Page is Not in Memory (Cache Miss) – Find a Free Frame

	frameID, isFromFreeList := b.getFrameID()
	if frameID == nil {
		return nil
	}


//if there arent any space in the frame so flush the page to the disk and then delete it from the frame to make space.
	if !isFromFreeList {
		currentPage := b.pages[*frameID]
		if currentPage != nil {
			if currentPage.isDirty {
				b.diskManager.WritePage(currentPage)
			}

			delete(b.pageTable, currentPage.id)
		}
	}

	page, err := b.diskManager.ReadPage(pageID)
	if err != nil {
		return nil
	}
	(*page).pinCount = 1
	b.pageTable[pageID] = *frameID
	b.pages[*frameID] = page

	return page
}

func (b *BufferPoolManager) UnpinPage(pageID PageID, isDirty bool) error {
	if frameID, ok := b.pageTable[pageID]; ok {
		page := b.pages[frameID]
		page.DecPinCount()

		if page.pinCount <= 0 {
			(*b.replacer).Unpin(frameID)
		}

		if page.isDirty || isDirty {
			page.isDirty = true
		} else {
			page.isDirty = false
		}

		return nil
	}

	return errors.New("Could not find page")
}

func (b *BufferPoolManager) FlushPage(pageID PageID) bool {
	if frameID, ok := b.pageTable[pageID]; ok {
		page := b.pages[frameID]
		page.DecPinCount()

		b.diskManager.WritePage(page)
		page.isDirty = false

		return true
	}

	return false
}

func (b *BufferPoolManager) NewPage() *Page {
	frameID, isFromFreeList := b.getFrameID()
	if frameID == nil {
		return nil
	}

	if !isFromFreeList {
		currentPage := b.pages[*frameID]
		if currentPage != nil {
			if currentPage.isDirty {
				b.diskManager.WritePage(currentPage)
			}

			delete(b.pageTable, currentPage.id)
		}
	}

	pageID := b.diskManager.AllocatePage()
	if pageID == nil {
		return nil
	}
	page := &Page{*pageID, 1, false, [pageSize]byte{}}

	b.pageTable[*pageID] = *frameID
	b.pages[*frameID] = page

	return page
}

func (b *BufferPoolManager) DeletePage(pageID PageID) error {
	var frameID FrameID
	var ok bool
	if frameID, ok = b.pageTable[pageID]; !ok {
		return nil
	}

	page := b.pages[frameID]

	if page.pinCount > 0 {
		return errors.New("Pin count greater than 0")
	}
	delete(b.pageTable, page.id)
	(*b.replacer).Pin(frameID)
	b.diskManager.DeallocatePage(pageID)

	b.freeList = append(b.freeList, frameID)

	return nil

}

func (b *BufferPoolManager) FlushAllpages() {
	for pageID := range b.pageTable {
		b.FlushPage(pageID)
	}
}

func (b *BufferPoolManager) getFrameID() (*FrameID, bool) {
	if len(b.freeList) > 0 {
		frameID, newFreeList := b.freeList[0], b.freeList[1:]
		b.freeList = newFreeList

		return &frameID, true
	}

	return (*b.replacer).Victim(), false
}

func NewBufferPoolManager(DiskManager DiskManager, clockReplacer *ClockReplacer) *BufferPoolManager {
	freeList := make([]FrameID, 0)
	pages := [MaxPoolSize]*Page{}
	for i := 0; i < MaxPoolSize; i++ {
		freeList = append(freeList, FrameID(i))
		pages[FrameID(i)] = nil
	}
	return &BufferPoolManager{DiskManager, pages, clockReplacer, freeList, make(map[PageID]FrameID)}
}