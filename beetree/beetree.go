package beetree

import (
	"fmt"
)

// CLRS B Trees
// Degree (t) is the min number of children a node can have.
// 2t is the max number of children a node can have.
// t-1 is the min number of keys a node can have.
// 2t-1 is the max number of keys a node can have.
// Root node can have min t-1 keya or 0 if no children.
// Intermediary nodes must have min t-1 keys and t children.
// Leaf nodes must have min t-1 keys.

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

func (n *Node) insertInSortedOrder(key Key) int {
	for i, k := range n.Keys {
		if key.K < k.K {
			rigthKeys := make([]Key, len(n.Keys[i:]))
			copy(rigthKeys, n.Keys[i:])

			n.Keys = n.Keys[:i]
			n.Keys = append(n.Keys, key)
			n.Keys = append(n.Keys, rigthKeys...)
			return i
		}
	}
	n.Keys = append(n.Keys, key)
	return len(n.Keys) - 1
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

	// Check if key already exists in current node.
	var keyExists bool
	var indexOfDuplicatedKey int
	for i, k := range node.Keys {
		if key.K == k.K {
			keyExists = true
			indexOfDuplicatedKey = i
			break
		}
	}

	if !keyExists {
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
			// left node otherwise in the right node.
			if key.K < middleKey.K {
				indexOfInsertedKey := node.insertInSortedOrder(key)
				if newSplitRigthChildNode != nil {
					// Insert the split child at the correct position in the left node
					insertPos := indexOfInsertedKey + 1
					if insertPos < len(node.Children) {
						node.Children = append(node.Children, nil)
						copy(node.Children[insertPos+1:], node.Children[insertPos:])
						node.Children[insertPos] = newSplitRigthChildNode
					} else {
						node.Children = append(node.Children, newSplitRigthChildNode)
					}
				}
			} else {
				indexOfInsertedKey := newRigthChildNode.insertInSortedOrder(key)
				if newSplitRigthChildNode != nil {
					// Insert the split child at the correct position in the right node by
					// using the index of the key that was inserted. The split child node should be
					// inserted one position after the index of the key.
					insertPos := indexOfInsertedKey + 1
					if insertPos < len(newRigthChildNode.Children) {
						newRigthChildNode.Children = append(newRigthChildNode.Children, nil)
						copy(newRigthChildNode.Children[insertPos+1:], newRigthChildNode.Children[insertPos:])
						newRigthChildNode.Children[insertPos] = newSplitRigthChildNode
					} else {
						newRigthChildNode.Children = append(newRigthChildNode.Children, newSplitRigthChildNode)
					}
				}
			}

			return newRigthChildNode, middleKey
		}
	}

	// When node has capacity to hold another key we just insert it.
	// We also check if one of its child nodes was split so that a new child node must
	// added to the list of children.
	if keyExists {
		node.Keys[indexOfDuplicatedKey] = key
	} else {
		indexOfInsertedKey := node.insertInSortedOrder(key)

		if newSplitRigthChildNode != nil {
			// Insert the new split child at the correct position
			// The new child should be inserted at indexOfInsertedKey + 1
			insertPos := indexOfInsertedKey + 1

			// Make room for the new child
			node.Children = append(node.Children, nil)
			copy(node.Children[insertPos+1:], node.Children[insertPos:])
			node.Children[insertPos] = newSplitRigthChildNode
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
	// Check if key exists in current node
	for i, k := range node.Keys {
		if key == k.K {
			return k
		}
		if key < k.K {
			// Key should be in left child
			if len(node.Children) > i {
				return bt.get(node.Children[i], key)
			}
			break
		}
	}

	// If we reach here and have children, key should be in rightmost child
	if len(node.Children) > 0 {
		return bt.get(node.Children[len(node.Keys)], key)
	}

	return Key{}
}

// PrintInLevelOrder prints the keys in the BeeTree in level order.
//
// Every printed key will have its parent index, node index and key value, all
// separated by colons.
//
// Example: 0:0:{20} -> 0[parent index]:0[node index]:{20}key
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

// find returns the index of the key if found in the current node or returns the index
// of the child node where key could be stored.
//
// If index of key is -1, the key was not found in the current node and index of child should be used
// to continue the finding.
func (bt *BeeTree) find(node *Node, key Key) (int, int) {
	for i, k := range node.Keys {
		if key.K == k.K {
			return i, -1
		}

		if key.K < k.K {
			return -1, i
		}
	}

	return -1, len(node.Keys)
}

// Delete deletes a key from the btree if found.
func (bt *BeeTree) Delete(key Key) {
	// If btree is empty, we return.
	if bt.Root == nil {
		return
	}

	bt.delete(bt.Root, key)
}

func (bt *BeeTree) delete(node *Node, key Key) {
	// Find if the key is in the current node or in which child node it could be.
	indexOfKey, indexOfChild := bt.find(node, key)

	// If found, we proceed with the deletion of the key.
	if indexOfKey >= 0 {
		// If node is a leaf node, it should not have children.
		// We just delete the key and return. Parent node should check if node is
		// underflow.
		if len(node.Children) == 0 {
			newKeys := make([]Key, 0, 2*bt.Degree-1)
			for i, k := range node.Keys {
				if i != indexOfKey {
					newKeys = append(newKeys, k)
				}
			}

			node.Keys = newKeys
			return
		}

		// If node has children, it is an intermediary node.
		// TODO: implement deletion for intermediary nodes.
		if len(node.Children) > 0 {
			return
		}
	}

	// If key is not in current node, we validate if the node has children otherwise this means that the key is not
	// in the tree.
	if len(node.Children) == 0 {
		return
	}

	// We move to the child where the key could be located. This index of child was returned from the find function.
	bt.delete(node.Children[indexOfChild], key)

	// Once returns, we check if child node is underflow due to the deletion of a key.
	// If not we return to finish the operation, otherwise if it is underflow, we redistribute or merge.
	if len(node.Children[indexOfChild].Keys) >= bt.Degree-1 {
		return
	}

	// Redistribution.
	// We find a left or rigth sibling node with enough keys so that we borrow one of their
	// keys that will be sent to the parent, and we take one from the parent for the underflow node.

	// We borrow from left sibling.
	// If this is not the first child.
	// If left sibling has enought keys.
	if indexOfChild > 0 && len(node.Children[indexOfChild-1].Keys) > bt.Degree-1 {
		// We get the key from the parent that will go to the underflow node.
		parentKey := node.Keys[indexOfChild-1]

		underflowNode := node.Children[indexOfChild]
		underflowNode.insertInSortedOrder(parentKey)

		leftSiblingNode := node.Children[indexOfChild-1]
		node.Keys[indexOfChild-1] = leftSiblingNode.Keys[len(leftSiblingNode.Keys)-1]

		leftSiblingNode.Keys = append(make([]Key, 0, 2*bt.Degree-1), leftSiblingNode.Keys[:len(leftSiblingNode.Keys)-1]...)

		return
	}

	// We use right sibling.
	// If this is not the last child.
	// If rigth sibling has enought keys.
	if indexOfChild < len(node.Keys) && len(node.Children[indexOfChild+1].Keys) > bt.Degree-1 {
		parentKey := node.Keys[indexOfChild]

		underflowNode := node.Children[indexOfChild]
		underflowNode.insertInSortedOrder(parentKey)

		rigthSiblingNode := node.Children[indexOfChild+1]
		node.Keys[indexOfChild] = rigthSiblingNode.Keys[0]

		rigthSiblingNode.Keys = append(make([]Key, 0, 2*bt.Degree-1), rigthSiblingNode.Keys[1:]...)
		return
	}

	// Merge
	// If left and rigth sibling node do not have enough keys to share, we must merge the current node with one of the siblings
	// and pull the separating key from the parent.
	indexOfKeyToPull := indexOfChild
	indexOfChild1 := indexOfChild
	indexOfChild2 := indexOfChild + 1

	// If this is the last child, we need to merge it with its
	// left sibling, since it does not have rigth sibling.
	if indexOfChild2 > len(node.Children) {
		indexOfKeyToPull = indexOfChild - 1
		indexOfChild1 = indexOfChild - 1
		indexOfChild2 = indexOfChild
	}

	// Create the new child node with current, sibling and parent key.
	mergedNode := NewNode(bt.Degree)
	mergedNode.insertInSortedOrder(node.Keys[indexOfKeyToPull])

	for _, k := range node.Children[indexOfChild1].Keys {
		mergedNode.insertInSortedOrder(k)
	}

	for _, k := range node.Children[indexOfChild2].Keys {
		mergedNode.insertInSortedOrder(k)
	}

	// Remove key from parent.
	newKeys := append(make([]Key, 0, bt.Degree-1), node.Keys[:indexOfKeyToPull]...)
	newKeys = append(newKeys, node.Keys[indexOfKeyToPull+1:]...)
	node.Keys = newKeys

	// Update child nodes.
	newChildren := make([]*Node, 0, 2*bt.Degree)
	for i, n := range node.Children {
		if i == indexOfChild1 {
			newChildren = append(newChildren, mergedNode)
		}

		if i < indexOfChild1 || i > indexOfChild2 {
			newChildren = append(newChildren, n)
		}
	}
	node.Children = newChildren
}
