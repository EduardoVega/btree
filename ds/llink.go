package ds

import "fmt"

type LList struct {
	rootNode *LListNode
	len      int
}

type LListNode struct {
	Data *Data
	Next *LListNode
}

func NewLList() LList {
	return LList{
		rootNode: nil,
		len:      0,
	}
}

func (l *LList) Len() int {
	return l.len
}

func (l *LList) Set(k int, v string) {
	if l.rootNode == nil {
		l.rootNode = &LListNode{
			Data: &Data{
				k,
				v,
			},
			Next: nil,
		}
	} else {
		node := l.rootNode
		for {
			if node.Next != nil {
				node = node.Next
			} else {
				newNode := LListNode{
					Data: &Data{
						k,
						v,
					},
				}

				node.Next = &newNode
				break
			}
		}
	}

	l.sort()
	l.len++
}

func (l *LList) Delete(k int) {
	if l.rootNode.Data.k == k {
		l.rootNode = l.rootNode.Next
		l.len--
		return
	}

	prevNode := l.rootNode
	node := l.rootNode.Next
	for {
		if node != nil {
			if node.Data.k == k {
				prevNode.Next = node.Next
				l.len--
				break
			}

			node = node.Next
			prevNode = prevNode.Next
		} else {
			break
		}
	}
}

func (l *LList) Get(k int) {
	node := l.rootNode
	indexCounter := 0

	for {
		if node != nil {
			if node.Data.k == k {
				fmt.Printf("%d: %d -> %s\n", indexCounter, node.Data.k, node.Data.v)
			}

			node = node.Next
			indexCounter++
		} else {
			break
		}
	}
}

func (l *LList) Print() {
	node := l.rootNode

	for {
		if node != nil {
			fmt.Printf("%d -> %s\n", node.Data.k, node.Data.v)
			node = node.Next
		} else {
			break
		}
	}
}

func (l *LList) sort() {

}
