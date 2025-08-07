package beetree

import (
	"math/rand"
	"testing"
)

const benchmarkTreeSize = 10000
const btreeDegree = 32

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []Key) {
	for _, v := range rand.Perm(n) {
		out = append(out, Key{K: v})
	}
	return
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		tr := NewBeetree(btreeDegree)
		for _, item := range insertP {
			tr.Insert(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

// func TestBTree(t *testing.T) {
// 	tr := NewBeetree(btreeDegree)
// 	const treeSize = 10000
// 	for i := 0; i < 10; i++ {
// 		if min := tr.Min(); min != nil {
// 			t.Fatalf("empty min, got %+v", min)
// 		}
// 		if max := tr.Max(); max != nil {
// 			t.Fatalf("empty max, got %+v", max)
// 		}
// 		for _, item := range perm(treeSize) {
// 			if x := tr.ReplaceOrInsert(item); x != nil {
// 				t.Fatal("insert found item", item)
// 			}
// 		}
// 		for _, item := range perm(treeSize) {
// 			if !tr.Has(item) {
// 				t.Fatal("has did not find item", item)
// 			}
// 		}
// 		for _, item := range perm(treeSize) {
// 			if x := tr.ReplaceOrInsert(item); x == nil {
// 				t.Fatal("insert didn't find item", item)
// 			}
// 		}
// 		if min, want := tr.Min(), Item(Int(0)); min != want {
// 			t.Fatalf("min: want %+v, got %+v", want, min)
// 		}
// 		if max, want := tr.Max(), Item(Int(treeSize-1)); max != want {
// 			t.Fatalf("max: want %+v, got %+v", want, max)
// 		}
// 		got := all(tr)
// 		want := rang(treeSize)
// 		if !reflect.DeepEqual(got, want) {
// 			t.Fatalf("mismatch:\n got: %v\nwant: %v", got, want)
// 		}

// 		gotrev := allrev(tr)
// 		wantrev := rangrev(treeSize)
// 		if !reflect.DeepEqual(gotrev, wantrev) {
// 			t.Fatalf("mismatch:\n got: %v\nwant: %v", got, want)
// 		}

// 		for _, item := range perm(treeSize) {
// 			if x := tr.Delete(item); x == nil {
// 				t.Fatalf("didn't find %v", item)
// 			}
// 		}
// 		if got = all(tr); len(got) > 0 {
// 			t.Fatalf("some left!: %v", got)
// 		}
// 	}
// }
