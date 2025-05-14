package main

import (
	"fmt"
	"go-ds/ds"
)

func main() {
	// myBTree := ds.NewBTree(2)
	// fmt.Println(myBTree)

	// myBTree.Insert(10, "v1")
	// fmt.Println(len(myBTree.RootNode.Data))

	linkedList := ds.NewLList()
	linkedList.Set(10, "v1")
	linkedList.Set(2, "v2")
	linkedList.Set(8, "v3")
	fmt.Println(linkedList.Len())
	linkedList.Print()
	linkedList.Get(10)

	linkedList.Delete(2)
	fmt.Println(linkedList.Len())
	linkedList.Print()

	linkedList.Get(2)
}
