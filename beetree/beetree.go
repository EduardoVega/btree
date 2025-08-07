package beetree

import (
	"fmt"
)

// CLRS B Trees
// Degree (t) is min number of children a node can have.
// 2t is max number of children a node can have.
// t-1 is min number of keys a node can have.
// 2t-1 is max number of keys a node can have.

type Key struct {
	K int
}

type Node struct {
	Keys     []Key
	Children []*Node
}

type BeeTree struct {
	Degree int
	Root   *Node
}

func NewNode(degree int) *Node {
	return &Node{
		make([]Key, 0, (2*degree)-1),
		make([]*Node, 0, 2*degree),
	}
}

func (n *Node) insertInSortedOrder(key Key) {
	for i, k := range n.Keys {
		if key.K < k.K {
			rigthKeys := make([]Key, len(n.Keys[i:]))
			copy(rigthKeys, n.Keys[i:])

			n.Keys = n.Keys[:i]
			n.Keys = append(n.Keys, key)
			n.Keys = append(n.Keys, rigthKeys...)
			return
		}
	}
	n.Keys = append(n.Keys, key)
}

func NewBeetree(degree int) *BeeTree {
	return &BeeTree{
		Degree: degree,
	}
}

func (bt *BeeTree) Insert(key Key) {
	if bt.Root == nil {
		bt.Root = NewNode(bt.Degree)
		bt.Root.Keys = append(bt.Root.Keys, key)
		return
	}

	newRigthChildNode, middleKey := bt.insert(bt.Root, key)
	// If a key has been returned to root, it means the tree has grown and a new
	// level must be created with a new root containing the returned key.
	if newRigthChildNode != nil {
		newRootNode := NewNode(bt.Degree)
		newRootNode.Keys = append(newRootNode.Keys, middleKey)
		newRootNode.Children = append(newRootNode.Children, bt.Root)
		newRootNode.Children = append(newRootNode.Children, newRigthChildNode)
		bt.Root = newRootNode
	}
}

func (bt *BeeTree) insert(node *Node, key Key) (*Node, Key) {
	// This holds the index of the child node that was split and it is used
	// to determine in what position to insert the new child node.
	var indexOfSplitNode = -1
	var newSplitRigthChildNode *Node

	// If node has children, we must traverse to find the node where the new key must be
	// inserted.
	if len(node.Children) > 0 {
		// We search first if key should be in left nodes.
		for i, k := range node.Keys {
			if key.K < k.K {
				indexOfSplitNode = i
				break
			}
		}

		// If no index, new key is bigger than existing keys so we move to most rigth node.
		// Most right node should always exist since child nodes == keys+1
		if indexOfSplitNode == -1 {
			indexOfSplitNode = len(node.Keys)
		}

		// If key and new child node is returned it means child node was split.
		// fmt.Printf("Index %d, len of nodes %d, keys of node %v, key to insert %d \n", indexOfSplitNode, len(node.Children), node.Keys, key)
		// bt.PrintInLevelOrder()
		// fmt.Println("------------------")

		newSplitRigthChildNode, key = bt.insert(node.Children[indexOfSplitNode], key)
		if newSplitRigthChildNode == nil {
			return nil, Key{}
		}
	}

	// If node is full (can not hold more keys), we must split it before adding the
	// new key. The split will get the middle key and create a new child node that
	// contains the keys bigger than the middle key. These will be returned to the parent
	// so that the middle can be inserted and new child node appended if it also has space
	// otherwise parent is also split.
	var newRigthChildNode *Node
	if len(node.Keys) == 2*bt.Degree-1 {
		// Store the middle key that needs to be sent upwards
		// to the parent node.
		middleIndex := bt.Degree - 1
		middleKey := node.Keys[middleIndex]

		// Create new child node with keys bigger than middle key and their children.
		newRigthChildNode = NewNode(bt.Degree)
		newRigthChildNode.Keys = append(newRigthChildNode.Keys, node.Keys[middleIndex+1:]...)
		if len(node.Children) >= middleIndex+1 {
			newRigthChildNode.Children = append(newRigthChildNode.Children, node.Children[middleIndex+1:]...)
		}

		// Set up existing node and update its keys to leave only smaller than middle key and their children.
		node.Keys = node.Keys[:middleIndex]
		if len(node.Children) >= middleIndex+1 {
			node.Children = node.Children[:middleIndex+1]
		}

		// Insert new key in left or right new child node.
		// If new key is less than the middle key it should be in the
		// left node otherwise in the rigth node.
		if key.K < middleKey.K {
			node.insertInSortedOrder(key)
			if newSplitRigthChildNode != nil {
				node.Children = append(node.Children, newSplitRigthChildNode)
			}
		} else {
			newRigthChildNode.insertInSortedOrder(key)
			if newSplitRigthChildNode != nil {
				newRigthChildNode.Children = append(newRigthChildNode.Children, newSplitRigthChildNode)
			}
		}

		return newRigthChildNode, middleKey
	}

	// When node has capacity to hold another key we just insert it.
	// We also check if one of its child nodes was split so that a new child node must
	// added to the list of children.
	node.insertInSortedOrder(key)
	if newSplitRigthChildNode != nil {
		// We check if the len of the children allows to shift children to the right
		// otherwise we just add the split node at the end of the slice
		if len(node.Children) > indexOfSplitNode+1 {
			lenRigthChildren := len(node.Children[indexOfSplitNode+1:])
			tempRigthChildren := make([]*Node, lenRigthChildren)
			copy(tempRigthChildren, node.Children[indexOfSplitNode+1:])

			node.Children = append(make([]*Node, 0, bt.Degree), node.Children[:indexOfSplitNode+1]...)
			node.Children = append(node.Children, newSplitRigthChildNode)
			node.Children = append(node.Children, tempRigthChildren...)
		} else {
			node.Children = append(node.Children, newSplitRigthChildNode)
		}
	}

	return nil, Key{}
}

func (bt *BeeTree) Get(key int) Key {
	if bt.Root == nil {
		return Key{}
	}

	return bt.get(bt.Root, key)
}

func (bt *BeeTree) get(node *Node, key int) Key {
	if len(node.Children) > 0 {
		// We search first if key should be in left nodes.
		for i, k := range node.Keys {
			if key < k.K {
				return bt.get(node.Children[i], key)
			}
		}

		// If new key is bigger to existing keys we move to most rigth node.
		return bt.get(node.Children[len(node.Keys)], key)
	}

	for _, k := range node.Keys {
		if key == k.K {
			return k
		}
	}

	return Key{}
}

func (bt *BeeTree) PrintInLevelOrder() {
	// Empty btree.
	if bt.Root == nil {
		return
	}

	// We create a slice with the nodes at each level, we start with root so
	// it is a slice of one node.
	nodes := make([]map[int]*Node, 0)
	nodes = append(nodes, map[int]*Node{-1: bt.Root})
	bt.printInLevelOrder(nodes)
}

func (bt *BeeTree) printInLevelOrder(nodes []map[int]*Node) {
	childrenNodes := make([]map[int]*Node, 0)

	// For every node in this level we print their keys and then create
	// a slice with the children nodes.
	for i, n := range nodes {
		for parentIndex, node := range n {
			for _, key := range node.Keys {
				fmt.Print(parentIndex, ":", i, ":", key, " ")
			}

			for _, c := range node.Children {
				childrenWithParentIndex := map[int]*Node{i: c}
				childrenNodes = append(childrenNodes, childrenWithParentIndex)
			}
		}
	}
	fmt.Println()

	if len(childrenNodes) > 0 {
		bt.printInLevelOrder(childrenNodes)
	}
}
