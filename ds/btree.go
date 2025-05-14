package ds

// Min num of keys = order-1
// Max num of keys = 2order-1
// Max num of children = 2order

type BTree struct {
	Order    int
	RootNode *Node
}

type Node struct {
	Data       []Data
	LeftChild  *Node
	RightChild *Node
}

func NewBTree(order int) *BTree {
	return &BTree{
		Order:    order,
		RootNode: nil,
	}
}

func (b *BTree) Insert(k int, v string) {
	if b.RootNode == nil {
		b.RootNode = &Node{
			Data:       make([]Data, 0, b.Order),
			LeftChild:  nil,
			RightChild: nil,
		}

		b.RootNode.Data = append(b.RootNode.Data, Data{k, v})
	}
}
