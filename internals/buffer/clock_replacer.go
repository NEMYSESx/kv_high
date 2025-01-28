package main

type FrameID int

type ClockReplacer struct {
	cList     *circularList
	clockHand **node
}