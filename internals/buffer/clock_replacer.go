package main

type FrameID int

type ClockReplacer struct {
	cList     *circularList
	clockHand **node
}

func (c *ClockReplacer) Victim() *FrameID {
	if c.cList.size == 0 {
		return nil
	}

	var victimFrameID *FrameID
	currentNode := (*c.clockHand)
	for {

		if currentNode.value.(bool) {
			currentNode.value = false
			c.clockHand = &currentNode.next
		} else {
			frameID := currentNode.key.(FrameID)
			victimFrameID = &frameID

			c.clockHand = &currentNode.next

			c.cList.remove(currentNode.key)
			return victimFrameID
		}
	}
}

func (c *ClockReplacer) Unpin(id FrameID) {
	if !c.cList.hasKey(id) {
		c.cList.insert(id, true)
		if c.cList.size == 1 {
			c.clockHand = &c.cList.head
		}
	}
}

func (c *ClockReplacer) Pin(id FrameID) {
	node := c.cList.find(id)
	if node == nil {
		return
	}

	if (*c.clockHand) == node {
		c.clockHand = &(*c.clockHand).next
	}
	c.cList.remove(id)

}

func (c *ClockReplacer) Size() int {
	return c.cList.size
}

func NewClockReplacer(poolSize int) *ClockReplacer {
	cList := newCircularList(poolSize)
	return &ClockReplacer{cList, &cList.head}
}