package main

import (
	"fmt"

	"go-ds/beetree"
	"go-ds/gbtree"
)

func main() {
	// BeeTree
	// startTime := time.Now()
	bt := beetree.NewBeetree(2)
	bt.Insert(beetree.Key{K: 50})
	bt.Insert(beetree.Key{K: 20})
	// bee.ReplaceOrInsert(beetree.BTreeItem{K: 20, V: "mule"})
	bt.Insert(beetree.Key{K: 80})
	bt.Insert(beetree.Key{K: 90})
	bt.Insert(beetree.Key{K: 70})
	bt.Insert(beetree.Key{K: 60})
	bt.Insert(beetree.Key{K: 65})
	bt.Insert(beetree.Key{K: 69})

	k := bt.Get(61)
	fmt.Println(k)

	// fmt.Println("Time:", time.Since(startTime))
	// fmt.Println(bee)
	bt.PrintInLevelOrder()

	fmt.Println("--------------------")

	// Google BTree
	// startTime = time.Now()
	gtree := gbtree.New(2)
	gtree.ReplaceOrInsert(gbtree.Int(50))
	gtree.ReplaceOrInsert(gbtree.Int(20))
	//gtree.ReplaceOrInsert(gbtree.Int(20))
	gtree.ReplaceOrInsert(gbtree.Int(80))
	gtree.ReplaceOrInsert(gbtree.Int(90))
	gtree.ReplaceOrInsert(gbtree.Int(70))
	gtree.ReplaceOrInsert(gbtree.Int(60))
	gtree.ReplaceOrInsert(gbtree.Int(65))
	gtree.ReplaceOrInsert(gbtree.Int(69))
	//gtree.ReplaceOrInsert(gbtree.Int(60))
	//gtree.ReplaceOrInsert(gbtree.Int(15))
	// fmt.Println("Time:", time.Since(startTime))
	// fmt.Println(gtree)
	gtree.LevelOrderTraversalPrint()

}
