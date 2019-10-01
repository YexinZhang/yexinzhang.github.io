package LinkList

import (
	"fmt"
)

var (
	nullMsg = "Usage: 'This Link list is a Null list.'"
)
type Linklist struct {
	Head *LinkNode
	Tail *LinkNode
	Size int
}

type LinkNode struct {
	Priv  *LinkNode
	Next  *LinkNode
	Value []byte
}

func NewLinkList() *Linklist {
	return &Linklist{}
}

func (this *Linklist) Null() bool {
	if this.Head == nil {
		fmt.Println(nullMsg)
		return true
	}
	return false
}

func (this *Linklist) Transfor() {
	if this.Null() {
		fmt.Println(nullMsg)
		return
	}
	tempNode := this.Head
	for tempNode.Next != nil {
		fmt.Println(tempNode.Value)
		tempNode = tempNode.Next
	}
	fmt.Println(tempNode.Value)
}



func (this *Linklist) Get(index int) *LinkNode {

	if this.Null() {
		fmt.Println(nullMsg)
		return nil
	}

	temp := this.Head
	for i := 0; i < index; i++ {
		temp = temp.Next
	}
	return temp
}


func (this *Linklist) InsertHead(value *LinkNode) {
	if this.Null(){
		this.Head = value
		this.Tail =value
		this.Size++
		return
	}
	this.Head.Priv = value
	value.Next = this.Head
	this.Head = value
	this.Size++
}

func (this *Linklist) InsertMid(index int, node *LinkNode) {

	if this.Null() {
		fmt.Println(nullMsg)
		this.Head = node
		this.Tail = node
		this.Size++
		return
	}
	temp := this.Get(index-1)
	node.Next = temp.Next
	node.Priv = temp
	temp.Next = node
	this.Size++
}

func (this *Linklist) InsertTail(node *LinkNode)  {
	this.Tail.Next = node
	node.Priv = this.Tail
	this.Tail = node
	this.Size++
}
