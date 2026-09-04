package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MaxLevel = 6

type Node struct {
	Key   int
	Value string
	Next  []*Node
}

type SkipList struct {
	head  *Node
	level int
}

func NewSkipList() *SkipList {
	return &SkipList{
		head: &Node{
			Next: make([]*Node, MaxLevel),
		},
		level: 1,
	}
}

func randomLevel() int {
	level := 1

	for level < MaxLevel && rand.Float64() < 0.5 {
		level++
	}

	return level
}

func (sl *SkipList) Insert(key int, value string) {
	update := make([]*Node, MaxLevel)
	current := sl.head

	// Find insertion position on every level
	for lvl := sl.level - 1; lvl >= 0; lvl-- {
		for current.Next[lvl] != nil &&
			current.Next[lvl].Key < key {
			current = current.Next[lvl]
		}

		update[lvl] = current
	}

	// Update existing key
	if next := current.Next[0]; next != nil && next.Key == key {
		next.Value = value
		return
	}

	nodeLevel := randomLevel()

	// Increase skip list height if needed
	if nodeLevel > sl.level {
		for i := sl.level; i < nodeLevel; i++ {
			update[i] = sl.head
		}
		sl.level = nodeLevel
	}

	newNode := &Node{
		Key:   key,
		Value: value,
		Next:  make([]*Node, nodeLevel),
	}

	// Insert node at all its levels
	for i := 0; i < nodeLevel; i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}
}

func (sl *SkipList) Search(key int) (string, bool) {
	current := sl.head

	for lvl := sl.level - 1; lvl >= 0; lvl-- {
		for current.Next[lvl] != nil &&
			current.Next[lvl].Key < key {
			current = current.Next[lvl]
		}
	}

	current = current.Next[0]

	if current != nil && current.Key == key {
		return current.Value, true
	}

	return "", false
}

func (sl *SkipList) PrintLevel0() {
	current := sl.head.Next[0]

	for current != nil {
		fmt.Printf("(%d,%s) -> ", current.Key, current.Value)
		current = current.Next[0]
	}

	fmt.Println("nil")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sl := NewSkipList()

	sl.Insert(10, "ten")
	sl.Insert(20, "twenty")
	sl.Insert(30, "thirty")
	sl.Insert(15, "fifteen")

	sl.PrintLevel0()

	if value, found := sl.Search(20); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}

	if value, found := sl.Search(100); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found")
	}
}
