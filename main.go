package main

import (
	"btree/beetree"
)

func main() {
	// BeeTree
	// startTime := time.Now()
	bt := beetree.NewBeetree(2)

	// 10, 20, 30, 40, 50, 60, 5, 15, 25, 35, 45, 55, 65, 75, 85, 95, 105
	bt.Insert(beetree.Key{K: 10})
	bt.Insert(beetree.Key{K: 20})
	bt.Insert(beetree.Key{K: 30})
	bt.Insert(beetree.Key{K: 40})
	bt.Insert(beetree.Key{K: 50})
	bt.Insert(beetree.Key{K: 60})
	bt.Insert(beetree.Key{K: 5})
	bt.Insert(beetree.Key{K: 15})
	bt.Insert(beetree.Key{K: 25})
	bt.Insert(beetree.Key{K: 35})
	bt.Insert(beetree.Key{K: 45})
	bt.Insert(beetree.Key{K: 55})
	bt.Insert(beetree.Key{K: 65})
	bt.Insert(beetree.Key{K: 75})
	bt.Insert(beetree.Key{K: 85})
	bt.Insert(beetree.Key{K: 95})
	bt.Insert(beetree.Key{K: 105})
	bt.PrintInLevelOrder()

	// Test left leaf deletion
	// Redistribution
	// bt.Delete(beetree.Key{K: 25})
	// bt.PrintInLevelOrder()

	// bt.Delete(beetree.Key{K: 35})
	// bt.PrintInLevelOrder()

	// bt.Delete(beetree.Key{K: 30})
	// bt.PrintInLevelOrder()

	// Test right leaf deletion
	// Redistribution
	// bt.Delete(beetree.Key{K: 15})
	// bt.PrintInLevelOrder()

	// bt.Delete(beetree.Key{K: 10})
	// bt.PrintInLevelOrder()

	// bt.Delete(beetree.Key{K: 5})
	// bt.PrintInLevelOrder()

	// Test leaf merge
	// bt.Delete(beetree.Key{K: 85})
	// bt.PrintInLevelOrder()

	// Test intermediate key deletion
	bt.Delete(beetree.Key{K: 20})
	bt.Delete(beetree.Key{K: 10})
	bt.Delete(beetree.Key{K: 15})
	bt.Delete(beetree.Key{K: 25})
	bt.Delete(beetree.Key{K: 5})

	bt.Delete(beetree.Key{K: 105})
	bt.Delete(beetree.Key{K: 95})
	bt.Delete(beetree.Key{K: 85})
	bt.Delete(beetree.Key{K: 75})

	bt.Delete(beetree.Key{K: 60})

	bt.PrintInLevelOrder()

	// bt.Insert(beetree.Key{K: 20})
	// bt.Insert(beetree.Key{K: 22})
	// bt.Insert(beetree.Key{K: 9})
	// bt.Insert(beetree.Key{K: 27})
	// bt.Insert(beetree.Key{K: 18})
	// bt.Insert(beetree.Key{K: 2})
	// bt.Insert(beetree.Key{K: 44})
	// bt.Insert(beetree.Key{K: 5})
	// bt.Insert(beetree.Key{K: 43})
	// bt.Insert(beetree.Key{K: 13})
	// bt.Insert(beetree.Key{K: 34})
	// bt.Insert(beetree.Key{K: 39})
	// bt.Insert(beetree.Key{K: 120})
	// bt.Insert(beetree.Key{K: 220})
	// bt.Insert(beetree.Key{K: 51})
	// bt.Insert(beetree.Key{K: 55})
	// bt.Insert(beetree.Key{K: 68})
	// bt.Insert(beetree.Key{K: 65})
	// bt.Insert(beetree.Key{K: 70})
	// bt.Insert(beetree.Key{K: 21})

	// bt.Insert(beetree.Key{K: 50})
	// bt.Insert(beetree.Key{K: 20})
	// // bee.ReplaceOrInsert(beetree.BTreeItem{K: 20, V: "mule"})
	// bt.Insert(beetree.Key{K: 80})
	// bt.Insert(beetree.Key{K: 90})
	// bt.Insert(beetree.Key{K: 70})
	// bt.Insert(beetree.Key{K: 60})
	// bt.Insert(beetree.Key{K: 65})
	// bt.Insert(beetree.Key{K: 69})

	//k := bt.Get(61)
	//fmt.Println(k)

	// fmt.Println("Time:", time.Since(startTime))
	// fmt.Println(bee)

	// fmt.Println("--------------------")

	// Google BTree
	// startTime = time.Now()
	// gtree := gbtree.New(2)
	//	gtree.ReplaceOrInsert(gbtree.Int(50))
	// gtree.ReplaceOrInsert(gbtree.Int(20))
	// //gtree.ReplaceOrInsert(gbtree.Int(20))
	// gtree.ReplaceOrInsert(gbtree.Int(80))
	// gtree.ReplaceOrInsert(gbtree.Int(90))
	// gtree.ReplaceOrInsert(gbtree.Int(70))
	// gtree.ReplaceOrInsert(gbtree.Int(60))
	// gtree.ReplaceOrInsert(gbtree.Int(65))
	// gtree.ReplaceOrInsert(gbtree.Int(69))

	// gtree.ReplaceOrInsert(gbtree.Int(44))
	// gtree.ReplaceOrInsert(gbtree.Int(4))
	// gtree.ReplaceOrInsert(gbtree.Int(28))
	// gtree.ReplaceOrInsert(gbtree.Int(3))
	// gtree.ReplaceOrInsert(gbtree.Int(15))
	// gtree.ReplaceOrInsert(gbtree.Int(30))
	// gtree.ReplaceOrInsert(gbtree.Int(48))
	// gtree.ReplaceOrInsert(gbtree.Int(17))
	// gtree.ReplaceOrInsert(gbtree.Int(38))
	// gtree.ReplaceOrInsert(gbtree.Int(23))

	// gtree.ReplaceOrInsert(gbtree.Int(21))
	// gtree.ReplaceOrInsert(gbtree.Int(20))
	// gtree.ReplaceOrInsert(gbtree.Int(22))
	// gtree.ReplaceOrInsert(gbtree.Int(9))
	// gtree.ReplaceOrInsert(gbtree.Int(27))
	// gtree.ReplaceOrInsert(gbtree.Int(18))
	// gtree.ReplaceOrInsert(gbtree.Int(2))
	// gtree.ReplaceOrInsert(gbtree.Int(44))
	// gtree.ReplaceOrInsert(gbtree.Int(5))
	// gtree.ReplaceOrInsert(gbtree.Int(43))
	// gtree.ReplaceOrInsert(gbtree.Int(13))
	// gtree.ReplaceOrInsert(gbtree.Int(34))
	//gtree.ReplaceOrInsert(gbtree.Int(39))

	//gtree.ReplaceOrInsert(gbtree.Int(60))
	//gtree.ReplaceOrInsert(gbtree.Int(15))
	// fmt.Println("Time:", time.Since(startTime))
	// fmt.Println(gtree)
	// gtree.LevelOrderTraversalPrint()

}
