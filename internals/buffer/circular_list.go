package main

type node struct {
	key   interface{}
	value interface{}
	next  *node
	prev  *node
}

type circularList struct {
	head     *node
	tail     *node
	size     int
	capacity int
}