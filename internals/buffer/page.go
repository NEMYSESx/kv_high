package main

type PageID int

const pageSize = 5

type Page struct {
	id       PageID
	pinCount int
	isDirty  bool
	data     [pageSize]byte
}

func (p *Page) PinCount() int {
	return p.pinCount
}

func (p *Page) ID() PageID {
	return p.id
}

func (p *Page) DecPinCount() {
	if p.pinCount > 0 {
		p.pinCount--
	}
}