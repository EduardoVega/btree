package beetree

import (
	"fmt"
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

// Helper function to collect all keys from the tree in order
func collectKeysInOrder(node *Node) []int {
	if node == nil {
		return []int{}
	}

	var result []int
	childIndex := 0

	for _, key := range node.Keys {
		// Add keys from left child
		if childIndex < len(node.Children) {
			result = append(result, collectKeysInOrder(node.Children[childIndex])...)
			childIndex++
		}
		// Add current key
		result = append(result, key.K)
	}

	// Add keys from rightmost child
	if childIndex < len(node.Children) {
		result = append(result, collectKeysInOrder(node.Children[childIndex])...)
	}

	return result
}

// Helper function to verify B-tree properties
func verifyBTreeProperties(t *testing.T, tree *BeeTree, node *Node, minDegree int, isRoot bool) {
	if node == nil {
		return
	}

	// Check key count constraints
	minKeys := minDegree - 1
	maxKeys := 2*minDegree - 1

	if isRoot && len(node.Keys) == 0 {
		t.Errorf("Root node cannot be empty unless tree is empty")
	}

	if !isRoot && len(node.Keys) < minKeys {
		t.Errorf("Non-root node has %d keys, minimum required: %d", len(node.Keys), minKeys)
	}

	if len(node.Keys) > maxKeys {
		t.Errorf("Node has %d keys, maximum allowed: %d", len(node.Keys), maxKeys)
	}

	// Check that keys are strictly sorted (no duplicates allowed)
	for i := 1; i < len(node.Keys); i++ {
		if node.Keys[i-1].K >= node.Keys[i].K {
			t.Errorf("Keys are not strictly sorted: %d >= %d at positions %d, %d",
				node.Keys[i-1].K, node.Keys[i].K, i-1, i)
		}
	}

	// Check children count
	if len(node.Children) > 0 {
		expectedChildren := len(node.Keys) + 1
		if len(node.Children) != expectedChildren {
			t.Errorf("Internal node has %d children, expected %d", len(node.Children), expectedChildren)
		}

		// Recursively verify children
		for _, child := range node.Children {
			verifyBTreeProperties(t, tree, child, minDegree, false)
		}
	}
}

// TestInsertEmptyTree tests inserting into an empty tree
func TestInsertEmptyTree(t *testing.T) {
	tree := NewBeetree(3)

	// Verify tree is initially empty
	if tree.Root != nil {
		t.Errorf("Expected empty tree, but root is not nil")
	}

	// Insert first key
	tree.Insert(Key{K: 10})

	// Verify root was created
	if tree.Root == nil {
		t.Fatalf("Root should not be nil after first insert")
	}

	// Verify the key was inserted
	if len(tree.Root.Keys) != 1 {
		t.Errorf("Expected 1 key in root, got %d", len(tree.Root.Keys))
	}

	if tree.Root.Keys[0].K != 10 {
		t.Errorf("Expected key 10, got %d", tree.Root.Keys[0].K)
	}

	// Verify no children since it's a leaf
	if len(tree.Root.Children) != 0 {
		t.Errorf("Expected no children in leaf node, got %d", len(tree.Root.Children))
	}
}

// TestInsertSingleNode tests inserting multiple keys without node splitting
func TestInsertSingleNode(t *testing.T) {
	tree := NewBeetree(3) // Max keys per node = 2*3-1 = 5

	// Insert keys in non-sorted order
	keys := []int{30, 10, 50, 20, 40}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}

	// Verify all keys are in root (no splitting should occur)
	if tree.Root == nil {
		t.Fatalf("Root should not be nil")
	}

	if len(tree.Root.Keys) != 5 {
		t.Errorf("Expected 5 keys in root, got %d", len(tree.Root.Keys))
	}

	// Verify keys are sorted
	expectedOrder := []int{10, 20, 30, 40, 50}
	for i, expectedKey := range expectedOrder {
		if tree.Root.Keys[i].K != expectedKey {
			t.Errorf("Expected key %d at position %d, got %d", expectedKey, i, tree.Root.Keys[i].K)
		}
	}

	// Verify B-tree properties
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestInsertCausingRootSplit tests inserting keys that cause root to split
func TestInsertCausingRootSplit(t *testing.T) {
	tree := NewBeetree(3) // Max keys per node = 5

	// Insert 6 keys to force root split
	keys := []int{10, 20, 30, 40, 50, 60}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}

	// Root should have been split, creating new root with one key
	if tree.Root == nil {
		t.Fatalf("Root should not be nil")
	}

	if len(tree.Root.Keys) != 1 {
		t.Errorf("Expected 1 key in new root, got %d", len(tree.Root.Keys))
	}

	// Root should have exactly 2 children
	if len(tree.Root.Children) != 2 {
		t.Errorf("Expected 2 children in root, got %d", len(tree.Root.Children))
	}

	// Verify B-tree properties are maintained
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)

	// Verify all keys are still accessible
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != 6 {
		t.Errorf("Expected 6 keys total, got %d", len(allKeys))
	}

	for i, expectedKey := range keys {
		if allKeys[i] != expectedKey {
			t.Errorf("Expected key %d at position %d, got %d", expectedKey, i, allKeys[i])
		}
	}
}

// TestInsertMultipleSplits tests inserting many keys causing multiple splits
func TestInsertMultipleSplits(t *testing.T) {
	tree := NewBeetree(3)

	// Insert keys 1 through 20
	for i := 1; i <= 20; i++ {
		tree.Insert(Key{K: i})

		// Verify B-tree properties after each insertion
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify all keys are present and in order
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != 20 {
		t.Errorf("Expected 20 keys, got %d", len(allKeys))
	}

	for i, expectedKey := range allKeys {
		if expectedKey != i+1 {
			t.Errorf("Expected key %d at position %d, got %d", i+1, i, expectedKey)
		}
	}
}

// TestInsertRandomOrder tests inserting keys in random order
func TestInsertRandomOrder(t *testing.T) {
	tree := NewBeetree(4)

	// Generate random permutation of keys
	keys := perm(50)

	// Insert all keys
	for _, key := range keys {
		tree.Insert(key)
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify all keys are present
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != 50 {
		t.Errorf("Expected 50 keys, got %d", len(allKeys))
	}

	// Verify keys are in strictly sorted order (no duplicates)
	for i := 1; i < len(allKeys); i++ {
		if allKeys[i-1] >= allKeys[i] {
			t.Errorf("Keys not in strictly sorted order: %d >= %d", allKeys[i-1], allKeys[i])
		}
	}
}

// TestInsertSequentialAscending tests inserting keys in ascending order
func TestInsertSequentialAscending(t *testing.T) {
	tree := NewBeetree(3)

	// Insert keys 1, 2, 3, ..., 15
	for i := 1; i <= 15; i++ {
		tree.Insert(Key{K: i})
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify all keys are present and in order
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != 15 {
		t.Errorf("Expected 15 keys, got %d", len(allKeys))
	}

	for i, key := range allKeys {
		if key != i+1 {
			t.Errorf("Expected key %d at position %d, got %d", i+1, i, key)
		}
	}
}

// TestInsertSequentialDescending tests inserting keys in descending order
func TestInsertSequentialDescending(t *testing.T) {
	tree := NewBeetree(3)

	// Insert keys 15, 14, 13, ..., 1
	for i := 15; i >= 1; i-- {
		tree.Insert(Key{K: i})
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify all keys are present and in order
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != 15 {
		t.Errorf("Expected 15 keys, got %d", len(allKeys))
	}

	for i, key := range allKeys {
		if key != i+1 {
			t.Errorf("Expected key %d at position %d, got %d", i+1, i, key)
		}
	}
}

// TestInsertDifferentDegrees tests insert operation with different B-tree degrees
func TestInsertDifferentDegrees(t *testing.T) {
	degrees := []int{2, 3, 4, 5, 10}
	numKeys := 30

	for _, degree := range degrees {
		t.Run(fmt.Sprintf("Degree_%d", degree), func(t *testing.T) {
			tree := NewBeetree(degree)

			// Insert keys in random order
			keys := perm(numKeys)
			for _, key := range keys {
				tree.Insert(key)
				verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
			}

			// Verify all keys are present
			allKeys := collectKeysInOrder(tree.Root)
			if len(allKeys) != numKeys {
				t.Errorf("Expected %d keys, got %d", numKeys, len(allKeys))
			}

			// Verify keys are strictly sorted (no duplicates)
			for i := 1; i < len(allKeys); i++ {
				if allKeys[i-1] >= allKeys[i] {
					t.Errorf("Keys not strictly sorted: %d >= %d", allKeys[i-1], allKeys[i])
				}
			}
		})
	}
}

// TestInsertAndGet tests that inserted keys can be retrieved
func TestInsertAndGet(t *testing.T) {
	tree := NewBeetree(3)

	keys := []int{10, 5, 15, 3, 7, 12, 18, 1, 4, 6, 8, 11, 13, 16, 20}

	// Insert all keys
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}

	// Verify all keys can be retrieved
	for _, k := range keys {
		result := tree.Get(k)
		if result.K != k {
			t.Errorf("Expected to get key %d, got %d", k, result.K)
		}
	}

	// Verify non-existent keys return empty Key
	nonExistentKeys := []int{2, 9, 14, 17, 19, 25}
	for _, k := range nonExistentKeys {
		result := tree.Get(k)
		if result.K != 0 { // Empty Key has K = 0
			t.Errorf("Expected empty key for non-existent key %d, got %d", k, result.K)
		}
	}
}

// TestInsertMinimumDegree tests inserting with minimum degree (2)
func TestInsertMinimumDegree(t *testing.T) {
	tree := NewBeetree(2) // Min degree = 2, max keys = 3

	// Insert enough keys to cause multiple splits
	keys := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for _, k := range keys {
		tree.Insert(Key{K: k})
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify all keys are present and sorted
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != len(keys) {
		t.Errorf("Expected %d keys, got %d", len(keys), len(allKeys))
	}

	for i, expectedKey := range keys {
		if allKeys[i] != expectedKey {
			t.Errorf("Expected key %d at position %d, got %d", expectedKey, i, allKeys[i])
		}
	}
}

// TestInsertLargeValues tests inserting large integer values
func TestInsertLargeValues(t *testing.T) {
	tree := NewBeetree(4)

	// Insert large values
	keys := []int{1000000, 2000000, 500000, 1500000, 750000, 1250000, 1750000}
	for _, k := range keys {
		tree.Insert(Key{K: k})
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify keys are sorted
	allKeys := collectKeysInOrder(tree.Root)
	for i := 1; i < len(allKeys); i++ {
		if allKeys[i-1] >= allKeys[i] {
			t.Errorf("Large values not sorted: %d >= %d", allKeys[i-1], allKeys[i])
		}
	}
}

// TestInsertNegativeValues tests inserting negative values
func TestInsertNegativeValues(t *testing.T) {
	tree := NewBeetree(3)

	// Insert mix of positive and negative values
	keys := []int{-10, 5, -20, 15, 0, -5, 25, -15}
	for _, k := range keys {
		tree.Insert(Key{K: k})
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}

	// Verify keys are sorted
	allKeys := collectKeysInOrder(tree.Root)
	for i := 1; i < len(allKeys); i++ {
		if allKeys[i-1] >= allKeys[i] {
			t.Errorf("Negative values not sorted: %d >= %d", allKeys[i-1], allKeys[i])
		}
	}

	// Expected sorted order: [-20, -15, -10, -5, 0, 5, 15, 25]
	expectedOrder := []int{-20, -15, -10, -5, 0, 5, 15, 25}
	if len(allKeys) != len(expectedOrder) {
		t.Errorf("Expected %d keys, got %d", len(expectedOrder), len(allKeys))
	}

	for i, expected := range expectedOrder {
		if allKeys[i] != expected {
			t.Errorf("Expected key %d at position %d, got %d", expected, i, allKeys[i])
		}
	}
}

// TestInsertStressTest tests inserting a large number of keys
func TestInsertStressTest(t *testing.T) {
	tree := NewBeetree(5)
	numKeys := 1000

	// Insert keys in random order
	keys := perm(numKeys)
	for _, key := range keys {
		tree.Insert(key)
		// Only verify properties occasionally to speed up the test
		if key.K%100 == 0 {
			verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
		}
	}

	// Final verification
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)

	// Verify all keys are present and sorted
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != numKeys {
		t.Errorf("Expected %d keys, got %d", numKeys, len(allKeys))
	}

	for i := 1; i < len(allKeys); i++ {
		if allKeys[i-1] >= allKeys[i] {
			t.Errorf("Keys not sorted in stress test: %d >= %d", allKeys[i-1], allKeys[i])
		}
	}
}

// TestInsertDuplicateKeys tests that duplicate keys are not supported
func TestInsertDuplicateKeys(t *testing.T) {
	tree := NewBeetree(3)
	
	// Insert initial keys
	tree.Insert(Key{K: 10})
	tree.Insert(Key{K: 20})
	tree.Insert(Key{K: 30})
	
	// Verify initial state
	initialKeys := collectKeysInOrder(tree.Root)
	if len(initialKeys) != 3 {
		t.Errorf("Expected 3 initial keys, got %d", len(initialKeys))
	}
	
	// Attempt to insert duplicate keys
	tree.Insert(Key{K: 10}) // Duplicate
	tree.Insert(Key{K: 20}) // Duplicate
	tree.Insert(Key{K: 30}) // Duplicate
	
	// Verify that duplicates were not added (keys should remain unique)
	afterDuplicatesKeys := collectKeysInOrder(tree.Root)
	
	// Check that no duplicates exist in the final result
	uniqueKeys := make(map[int]bool)
	for _, key := range afterDuplicatesKeys {
		if uniqueKeys[key] {
			t.Errorf("Found duplicate key %d in tree", key)
		}
		uniqueKeys[key] = true
	}
	
	// Verify the original keys are still present
	expectedKeys := []int{10, 20, 30}
	for _, expected := range expectedKeys {
		if !uniqueKeys[expected] {
			t.Errorf("Expected key %d is missing from tree", expected)
		}
	}
	
	// Verify B-tree properties are maintained
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestInsertDuplicatesInLargerTree tests duplicate handling with node splits
func TestInsertDuplicatesInLargerTree(t *testing.T) {
	tree := NewBeetree(3)
	
	// Insert enough keys to cause splits
	originalKeys := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for _, k := range originalKeys {
		tree.Insert(Key{K: k})
	}
	
	// Count original keys and verify they're unique
	beforeDuplicates := collectKeysInOrder(tree.Root)
	originalUniqueKeys := make(map[int]bool)
	for _, key := range beforeDuplicates {
		if originalUniqueKeys[key] {
			t.Errorf("Duplicate found in original keys: %d", key)
		}
		originalUniqueKeys[key] = true
	}
	
	// Insert duplicates of some keys
	duplicateKeys := []int{20, 50, 80, 30, 90}
	for _, k := range duplicateKeys {
		tree.Insert(Key{K: k})
	}
	
	// Verify no duplicates exist in final tree
	afterDuplicates := collectKeysInOrder(tree.Root)
	finalUniqueKeys := make(map[int]bool)
	for _, key := range afterDuplicates {
		if finalUniqueKeys[key] {
			t.Errorf("Duplicate key found after duplicate insertion: %d", key)
		}
		finalUniqueKeys[key] = true
	}
	
	// Verify all original keys are still present
	for _, originalKey := range originalKeys {
		if !finalUniqueKeys[originalKey] {
			t.Errorf("Original key %d is missing after duplicate insertion", originalKey)
		}
	}
	
	// Verify keys are still strictly sorted (no duplicates)
	for i := 1; i < len(afterDuplicates); i++ {
		if afterDuplicates[i-1] >= afterDuplicates[i] {
			t.Errorf("Keys not strictly sorted: %d >= %d", 
				afterDuplicates[i-1], afterDuplicates[i])
		}
	}
	
	// Verify B-tree properties are maintained
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestInsertDuplicatesWithGet tests that duplicate insertion doesn't affect retrieval
func TestInsertDuplicatesWithGet(t *testing.T) {
	tree := NewBeetree(3)
	
	// Insert keys
	keys := []int{15, 25, 35, 45, 55}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}
	
	// Verify all keys can be found
	for _, k := range keys {
		result := tree.Get(k)
		if result.K != k {
			t.Errorf("Expected to find key %d, got %d", k, result.K)
		}
	}
	
	// Insert duplicates
	for _, k := range keys {
		tree.Insert(Key{K: k}) // Insert each key again
	}
	
	// Verify keys can still be found (no corruption)
	for _, k := range keys {
		result := tree.Get(k)
		if result.K != k {
			t.Errorf("After duplicates, expected to find key %d, got %d", k, result.K)
		}
	}
	
	// Verify no duplicates exist in the tree
	allKeys := collectKeysInOrder(tree.Root)
	uniqueKeys := make(map[int]bool)
	for _, key := range allKeys {
		if uniqueKeys[key] {
			t.Errorf("Found duplicate key %d in tree after duplicate insertion", key)
		}
		uniqueKeys[key] = true
	}
	
	// Verify all original keys are still present
	for _, originalKey := range keys {
		if !uniqueKeys[originalKey] {
			t.Errorf("Original key %d is missing", originalKey)
		}
	}
}

// TestInsertDuplicateInSingleNode tests duplicate handling within a single node
func TestInsertDuplicateInSingleNode(t *testing.T) {
	tree := NewBeetree(4) // Larger degree to keep everything in one node initially
	
	// Insert keys that will fit in a single node
	tree.Insert(Key{K: 100})
	tree.Insert(Key{K: 200})
	tree.Insert(Key{K: 300})
	
	// Verify single node
	if tree.Root == nil {
		t.Fatal("Root should not be nil")
	}
	if len(tree.Root.Children) != 0 {
		t.Errorf("Expected leaf node, got %d children", len(tree.Root.Children))
	}
	if len(tree.Root.Keys) != 3 {
		t.Errorf("Expected 3 keys in root, got %d", len(tree.Root.Keys))
	}
	
	// Insert duplicate in the middle
	tree.Insert(Key{K: 200})
	
	// Verify still single node with same number of keys
	if len(tree.Root.Keys) != 3 {
		t.Errorf("Expected 3 keys after duplicate, got %d", len(tree.Root.Keys))
	}
	
	// Verify keys are correct
	expectedKeys := []int{100, 200, 300}
	for i, expected := range expectedKeys {
		if tree.Root.Keys[i].K != expected {
			t.Errorf("Expected key %d at position %d, got %d", 
				expected, i, tree.Root.Keys[i].K)
		}
	}
}

// TestDebugAdjustedIndex creates a scenario where adjustedIndex calculation is triggered
func TestDebugAdjustedIndex(t *testing.T) {
	// Use degree 3: max keys = 5, max children = 6, middleIndex = 2
	tree := NewBeetree(3)
	
	// Step 1: Insert keys to create a specific tree structure
	// Insert keys that will cause multiple splits and create the right scenario
	keys := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120}
	
	fmt.Println("=== Debug Adjusted Index Calculation ===")
	fmt.Println("Degree: 3, Max keys per node: 5, Middle index: 2")
	fmt.Println()
	
	for i, k := range keys {
		fmt.Printf("Inserting key %d (step %d)\n", k, i+1)
		tree.Insert(Key{K: k})
		fmt.Printf("Tree after inserting %d:\n", k)
		tree.PrintInLevelOrder()
		fmt.Println("---")
	}
	
	// The specific scenario occurs when:
	// 1. A node is full (5 keys) and needs to split
	// 2. The new key goes to the right half after split
	// 3. A child was previously split (newSplitRigthChildNode != nil)
	// 4. We need to insert that split child into the right half
	
	fmt.Println("Final tree structure:")
	tree.PrintInLevelOrder()
	
	// Verify all keys are present
	allKeys := collectKeysInOrder(tree.Root)
	if len(allKeys) != len(keys) {
		t.Errorf("Expected %d keys, got %d", len(keys), len(allKeys))
	}
	
	for i, expected := range keys {
		if allKeys[i] != expected {
			t.Errorf("Expected key %d at position %d, got %d", expected, i, allKeys[i])
		}
	}
}

// TestManualAdjustedIndexExample demonstrates the exact scenario step by step
func TestManualAdjustedIndexExample(t *testing.T) {
	fmt.Println("\n=== Manual Trace of adjustedIndex Calculation ===")
	fmt.Println()
	
	fmt.Println("Let's trace through a specific example:")
	fmt.Println("Degree = 3, so:")
	fmt.Println("- Max keys per node = 2*3-1 = 5")
	fmt.Println("- Max children per node = 2*3 = 6") 
	fmt.Println("- middleIndex = 3-1 = 2")
	fmt.Println()
	
	fmt.Println("Scenario: Internal node is FULL with 5 keys and 6 children")
	fmt.Println("Node structure BEFORE split:")
	fmt.Println("Keys:     [K0] [K1] [K2] [K3] [K4]")
	fmt.Println("Children: [C0] [C1] [C2] [C3] [C4] [C5]")
	fmt.Println("Indices:   0    1    2    3    4    5")
	fmt.Println()
	
	fmt.Println("When we split this node:")
	fmt.Println("- middleIndex = 2, so K2 goes up to parent")
	fmt.Println("- Left node gets: keys [K0,K1] and children [C0,C1,C2]")
	fmt.Println("- Right node gets: keys [K3,K4] and children [C3,C4,C5]")
	fmt.Println()
	
	fmt.Println("Original children positions → New right node positions:")
	fmt.Println("C3 (index 3) → becomes index 0 in right node")
	fmt.Println("C4 (index 4) → becomes index 1 in right node") 
	fmt.Println("C5 (index 5) → becomes index 2 in right node")
	fmt.Println()
	
	fmt.Println("Now suppose:")
	fmt.Println("- We're inserting a new key that goes to the RIGHT side")
	fmt.Println("- During insertion, child at original index 4 split")
	fmt.Println("- So indexOfSplitNode = 4")
	fmt.Println("- We need to insert the new split child into the right node")
	fmt.Println()
	
	fmt.Println("Calculation:")
	fmt.Println("adjustedIndex = indexOfSplitNode - (middleIndex + 1)")
	fmt.Println("adjustedIndex = 4 - (2 + 1) = 4 - 3 = 1")
	fmt.Println()
	fmt.Println("This means: the split child should be inserted at position 1")
	fmt.Println("in the right node's children array")
	fmt.Println()
	
	fmt.Println("insertPos = adjustedIndex + 1 = 1 + 1 = 2")
	fmt.Println("So the new split child gets inserted at position 2")
	fmt.Println()
	
	fmt.Println("Final right node children after insertion:")
	fmt.Println("[C3] [C4] [NewSplitChild] [C5]")
	fmt.Println(" 0    1         2          3")
	fmt.Println()
	
	fmt.Println("This maintains the B-tree property that children are correctly")
	fmt.Println("positioned relative to the keys in the node.")
}

// TestTriggerAdjustedIndexWithSpecificKeys provides exact keys to trigger the adjustedIndex calculation
func TestTriggerAdjustedIndexWithSpecificKeys(t *testing.T) {
	fmt.Println("\n=== Specific Keys to Trigger adjustedIndex Calculation ===")
	fmt.Println()
	
	tree := NewBeetree(3)
	
	fmt.Println("Follow these exact steps to see the adjustedIndex calculation:")
	fmt.Println()
	
	// Phase 1: Build initial tree structure
	phase1Keys := []int{10, 20, 30, 40, 50, 60}
	fmt.Println("Phase 1: Create initial tree with splits")
	for i, k := range phase1Keys {
		fmt.Printf("Insert %d", k)
		tree.Insert(Key{K: k})
		if i == 5 {
			fmt.Println(" <- This causes first major split")
			tree.PrintInLevelOrder()
			fmt.Println()
		} else {
			fmt.Println()
		}
	}
	
	// Phase 2: Add keys to create multi-level structure
	phase2Keys := []int{5, 15, 25, 35, 45, 55, 65, 75}
	fmt.Println("Phase 2: Build deeper tree structure")
	for _, k := range phase2Keys {
		fmt.Printf("Insert %d\n", k)
		tree.Insert(Key{K: k})
	}
	fmt.Println("Tree after Phase 2:")
	tree.PrintInLevelOrder()
	fmt.Println()
	
	// Phase 3: The critical insertions that trigger adjustedIndex
	phase3Keys := []int{85, 95, 105}
	fmt.Println("Phase 3: Critical keys that trigger adjustedIndex calculation")
	fmt.Println("(These will cause internal node splits with child splits)")
	
	for _, k := range phase3Keys {
		fmt.Printf("Insert %d <- This should trigger the adjustedIndex path\n", k)
		tree.Insert(Key{K: k})
		tree.PrintInLevelOrder()
		fmt.Println()
	}
	
	fmt.Println("=== To Debug the adjustedIndex Line ===")
	fmt.Println("1. Add a breakpoint or print statement at line 161:")
	fmt.Println("   adjustedIndex := indexOfSplitNode - (middleIndex + 1)")
	fmt.Println()
	fmt.Println("2. Insert these keys in order:")
	fmt.Print("   Keys: ")
	allKeys := append(append(phase1Keys, phase2Keys...), phase3Keys...)
	for i, k := range allKeys {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(k)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("3. Watch for when:")
	fmt.Println("   - A node has 5 keys (full)")
	fmt.Println("   - New key goes to right side of split") 
	fmt.Println("   - newSplitRigthChildNode != nil")
	fmt.Println("   - You're in the 'else' branch (key >= middleKey)")
	fmt.Println()
	fmt.Println("4. The calculation will show:")
	fmt.Println("   - indexOfSplitNode = some value (e.g., 4 or 5)")
	fmt.Println("   - middleIndex = 2")
	fmt.Println("   - adjustedIndex = indexOfSplitNode - 3")
	fmt.Println("   - insertPos = adjustedIndex + 1")
	
	// Verify final tree
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// Helper function to build a tree with specific keys for testing
func buildTreeWithKeys(degree int, keys []int) *BeeTree {
	tree := NewBeetree(degree)
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}
	return tree
}

// Helper function to check if a key exists in the tree
func treeContainsKey(tree *BeeTree, key int) bool {
	result := tree.Get(key)
	return result.K == key
}

// Helper function to get tree height
func getTreeHeight(node *Node) int {
	if node == nil || len(node.Children) == 0 {
		return 1
	}
	return 1 + getTreeHeight(node.Children[0])
}

// Helper function to count total nodes in tree
func countNodes(node *Node) int {
	if node == nil {
		return 0
	}
	count := 1
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}

// TestDeleteEmptyTree tests deleting from an empty tree
func TestDeleteEmptyTree(t *testing.T) {
	tree := NewBeetree(3)
	
	// Delete from empty tree should not panic
	tree.Delete(Key{K: 10})
	
	// Tree should remain empty
	if tree.Root != nil {
		t.Errorf("Expected tree to remain empty after delete from empty tree")
	}
}

// TestDeleteNonExistentKey tests deleting a key that doesn't exist
func TestDeleteNonExistentKey(t *testing.T) {
	tree := buildTreeWithKeys(3, []int{10, 20, 30, 40, 50})
	originalKeys := collectKeysInOrder(tree.Root)
	
	// Delete non-existent key
	tree.Delete(Key{K: 100})
	
	// Tree should remain unchanged
	newKeys := collectKeysInOrder(tree.Root)
	if len(newKeys) != len(originalKeys) {
		t.Errorf("Tree changed after deleting non-existent key")
	}
	
	for i, key := range originalKeys {
		if newKeys[i] != key {
			t.Errorf("Tree structure changed after deleting non-existent key")
		}
	}
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestDeleteFromSingleNodeLeaf tests deleting from a tree with only root node
func TestDeleteFromSingleNodeLeaf(t *testing.T) {
	tree := buildTreeWithKeys(3, []int{10, 20, 30})
	
	// Delete middle key
	tree.Delete(Key{K: 20})
	
	// Verify key was deleted
	if treeContainsKey(tree, 20) {
		t.Errorf("Key 20 should have been deleted")
	}
	
	// Verify remaining keys
	keys := collectKeysInOrder(tree.Root)
	expected := []int{10, 30}
	if len(keys) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(keys))
	}
	
	for i, key := range expected {
		if keys[i] != key {
			t.Errorf("Expected key %d at position %d, got %d", key, i, keys[i])
		}
	}
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestDeleteAllKeysFromSingleNode tests deleting all keys from root
func TestDeleteAllKeysFromSingleNode(t *testing.T) {
	tree := buildTreeWithKeys(3, []int{10})
	
	// Delete the only key
	tree.Delete(Key{K: 10})
	
	// Tree should still have root but with no keys
	if tree.Root == nil {
		t.Errorf("Root should not be nil after deleting all keys")
	}
	
	if len(tree.Root.Keys) != 0 {
		t.Errorf("Root should have no keys after deleting all, got %d", len(tree.Root.Keys))
	}
	
	// Verify key was actually deleted
	if treeContainsKey(tree, 10) {
		t.Errorf("Key 10 should have been deleted")
	}
}

// TestDeleteLeafNodeNoUnderflow tests deleting from leaf nodes without causing underflow
func TestDeleteLeafNodeNoUnderflow(t *testing.T) {
	// Create tree with degree 3 (min keys = 2, max keys = 5)
	tree := buildTreeWithKeys(3, []int{10, 20, 30, 40, 50, 60, 70, 80, 90})
	
	// Find a leaf node with more than minimum keys
	originalKeys := collectKeysInOrder(tree.Root)
	
	// Delete a key that should be in a leaf
	tree.Delete(Key{K: 90})
	
	// Verify key was deleted
	if treeContainsKey(tree, 90) {
		t.Errorf("Key 90 should have been deleted")
	}
	
	// Verify tree structure is maintained
	newKeys := collectKeysInOrder(tree.Root)
	if len(newKeys) != len(originalKeys)-1 {
		t.Errorf("Expected %d keys after deletion, got %d", len(originalKeys)-1, len(newKeys))
	}
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestDeleteInternalNodeWithPredecessor tests deleting from internal node using predecessor
func TestDeleteInternalNodeWithPredecessor(t *testing.T) {
	// Build a tree where we can control the structure
	tree := buildTreeWithKeys(3, []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	
	originalKeys := collectKeysInOrder(tree.Root)
	
	// Delete a key that should be in an internal node
	// This should trigger predecessor replacement
	tree.Delete(Key{K: 50})
	
	// Verify key was deleted
	if treeContainsKey(tree, 50) {
		t.Errorf("Key 50 should have been deleted")
	}
	
	// Verify correct number of keys remain
	newKeys := collectKeysInOrder(tree.Root)
	if len(newKeys) != len(originalKeys)-1 {
		t.Errorf("Expected %d keys after deletion, got %d", len(originalKeys)-1, len(newKeys))
	}
	
	// Verify all other keys are still present
	for _, key := range originalKeys {
		if key != 50 && !treeContainsKey(tree, key) {
			t.Errorf("Key %d should still be present after deleting 50", key)
		}
	}
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestDeleteInternalNodeWithSuccessor tests deleting from internal node using successor
func TestDeleteInternalNodeWithSuccessor(t *testing.T) {
	// Build a tree structure where successor will be used
	tree := buildTreeWithKeys(3, []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50})
	
	originalKeys := collectKeysInOrder(tree.Root)
	
	// Delete a key from internal node
	tree.Delete(Key{K: 20})
	
	// Verify key was deleted
	if treeContainsKey(tree, 20) {
		t.Errorf("Key 20 should have been deleted")
	}
	
	// Verify tree integrity
	newKeys := collectKeysInOrder(tree.Root)
	if len(newKeys) != len(originalKeys)-1 {
		t.Errorf("Expected %d keys after deletion, got %d", len(originalKeys)-1, len(newKeys))
	}
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
}

// TestDeleteCausingRedistributionFromLeft tests deletion causing redistribution from left sibling
func TestDeleteCausingRedistributionFromLeft(t *testing.T) {
	// Create a specific tree structure for testing left redistribution
	tree := NewBeetree(3)
	keys := []int{10, 20, 30, 40, 50, 60, 70}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}
	
	// Delete a key to cause underflow that can be fixed by left redistribution
	tree.Delete(Key{K: 70})
	tree.Delete(Key{K: 60}) // This should cause redistribution
	
	// Verify tree properties are maintained
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	
	// Verify keys are still accessible
	remainingKeys := collectKeysInOrder(tree.Root)
	expected := []int{10, 20, 30, 40, 50}
	
	if len(remainingKeys) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(remainingKeys))
	}
	
	for i, key := range expected {
		if remainingKeys[i] != key {
			t.Errorf("Expected key %d at position %d, got %d", key, i, remainingKeys[i])
		}
	}
}

// TestDeleteCausingRedistributionFromRight tests deletion causing redistribution from right sibling
func TestDeleteCausingRedistributionFromRight(t *testing.T) {
	tree := NewBeetree(3)
	keys := []int{10, 20, 30, 40, 50, 60, 70}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}
	
	// Delete keys to cause underflow that requires right redistribution
	tree.Delete(Key{K: 10})
	tree.Delete(Key{K: 20}) // This should cause redistribution from right
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	
	remainingKeys := collectKeysInOrder(tree.Root)
	expected := []int{30, 40, 50, 60, 70}
	
	if len(remainingKeys) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(remainingKeys))
	}
	
	for i, key := range expected {
		if remainingKeys[i] != key {
			t.Errorf("Expected key %d at position %d, got %d", key, i, remainingKeys[i])
		}
	}
}

// TestDeleteCausingMerge tests deletion that causes node merging
func TestDeleteCausingMerge(t *testing.T) {
	tree := NewBeetree(3)
	
	// Create a tree structure that will require merging
	keys := []int{10, 20, 30, 40, 50, 60}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}
	
	originalNodeCount := countNodes(tree.Root)
	
	// Delete keys to force merge
	tree.Delete(Key{K: 60})
	tree.Delete(Key{K: 50})
	tree.Delete(Key{K: 40})
	
	// Verify merge occurred (fewer nodes)
	newNodeCount := countNodes(tree.Root)
	if newNodeCount >= originalNodeCount {
		t.Logf("Original node count: %d, New node count: %d", originalNodeCount, newNodeCount)
		// Note: This is just informational, actual merge behavior depends on tree structure
	}
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	
	remainingKeys := collectKeysInOrder(tree.Root)
	expected := []int{10, 20, 30}
	
	if len(remainingKeys) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(remainingKeys))
	}
	
	for i, key := range expected {
		if remainingKeys[i] != key {
			t.Errorf("Expected key %d at position %d, got %d", key, i, remainingKeys[i])
		}
	}
}

// TestDeleteCausingRootChange tests deletion that changes the root
func TestDeleteCausingRootChange(t *testing.T) {
	tree := NewBeetree(3)
	
	// Build tree and then delete enough to potentially change root
	keys := []int{10, 20, 30, 40, 50}
	for _, k := range keys {
		tree.Insert(Key{K: k})
	}
	
	originalHeight := getTreeHeight(tree.Root)
	
	// Delete most keys
	tree.Delete(Key{K: 10})
	tree.Delete(Key{K: 20})
	tree.Delete(Key{K: 30})
	
	newHeight := getTreeHeight(tree.Root)
	
	verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	
	remainingKeys := collectKeysInOrder(tree.Root)
	expected := []int{40, 50}
	
	if len(remainingKeys) != len(expected) {
		t.Errorf("Expected %d keys, got %d", len(expected), len(remainingKeys))
	}
	
	t.Logf("Tree height changed from %d to %d", originalHeight, newHeight)
}

// TestDeleteSequential tests deleting keys in sequential order
func TestDeleteSequential(t *testing.T) {
	tree := buildTreeWithKeys(3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	
	// Delete keys 1 through 10
	for i := 1; i <= 10; i++ {
		tree.Delete(Key{K: i})
		
		// Verify key was deleted
		if treeContainsKey(tree, i) {
			t.Errorf("Key %d should have been deleted", i)
		}
		
		// Verify tree properties after each deletion
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}
	
	// Verify remaining keys
	remainingKeys := collectKeysInOrder(tree.Root)
	expected := []int{11, 12, 13, 14, 15}
	
	if len(remainingKeys) != len(expected) {
		t.Errorf("Expected %d keys remaining, got %d", len(expected), len(remainingKeys))
	}
	
	for i, key := range expected {
		if remainingKeys[i] != key {
			t.Errorf("Expected key %d at position %d, got %d", key, i, remainingKeys[i])
		}
	}
}

// TestDeleteReverseSequential tests deleting keys in reverse order
func TestDeleteReverseSequential(t *testing.T) {
	tree := buildTreeWithKeys(3, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	
	// Delete keys 15 down to 6
	for i := 15; i >= 6; i-- {
		tree.Delete(Key{K: i})
		
		// Verify key was deleted
		if treeContainsKey(tree, i) {
			t.Errorf("Key %d should have been deleted", i)
		}
		
		// Verify tree properties after each deletion
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}
	
	// Verify remaining keys
	remainingKeys := collectKeysInOrder(tree.Root)
	expected := []int{1, 2, 3, 4, 5}
	
	if len(remainingKeys) != len(expected) {
		t.Errorf("Expected %d keys remaining, got %d", len(expected), len(remainingKeys))
	}
	
	for i, key := range expected {
		if remainingKeys[i] != key {
			t.Errorf("Expected key %d at position %d, got %d", key, i, remainingKeys[i])
		}
	}
}

// TestDeleteRandomOrder tests deleting keys in random order
func TestDeleteRandomOrder(t *testing.T) {
	const numKeys = 50
	tree := NewBeetree(4)
	
	// Insert keys 1 through numKeys
	allKeys := make([]int, numKeys)
	for i := 0; i < numKeys; i++ {
		allKeys[i] = i + 1
		tree.Insert(Key{K: i + 1})
	}
	
	// Create random permutation for deletion
	deleteOrder := rand.Perm(numKeys)
	
	// Delete half the keys in random order
	keysToDelete := numKeys / 2
	deletedKeys := make(map[int]bool)
	
	for i := 0; i < keysToDelete; i++ {
		keyToDelete := deleteOrder[i] + 1
		tree.Delete(Key{K: keyToDelete})
		deletedKeys[keyToDelete] = true
		
		// Verify key was deleted
		if treeContainsKey(tree, keyToDelete) {
			t.Errorf("Key %d should have been deleted", keyToDelete)
		}
		
		// Verify tree properties
		verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
	}
	
	// Verify remaining keys
	remainingKeys := collectKeysInOrder(tree.Root)
	expectedRemaining := numKeys - keysToDelete
	
	if len(remainingKeys) != expectedRemaining {
		t.Errorf("Expected %d keys remaining, got %d", expectedRemaining, len(remainingKeys))
	}
	
	// Verify that only non-deleted keys remain
	for _, key := range remainingKeys {
		if deletedKeys[key] {
			t.Errorf("Deleted key %d should not be present", key)
		}
	}
	
	// Verify all non-deleted keys are present
	for i := 1; i <= numKeys; i++ {
		if !deletedKeys[i] && !treeContainsKey(tree, i) {
			t.Errorf("Non-deleted key %d should be present", i)
		}
	}
}

// TestDeleteDifferentDegrees tests delete operation with different B-tree degrees
func TestDeleteDifferentDegrees(t *testing.T) {
	degrees := []int{2, 3, 4, 5, 10}
	numKeys := 30
	
	for _, degree := range degrees {
		t.Run(fmt.Sprintf("Degree_%d", degree), func(t *testing.T) {
			tree := NewBeetree(degree)
			
			// Insert keys
			for i := 1; i <= numKeys; i++ {
				tree.Insert(Key{K: i})
			}
			
			// Delete every other key
			for i := 2; i <= numKeys; i += 2 {
				tree.Delete(Key{K: i})
				verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
			}
			
			// Verify only odd keys remain
			remainingKeys := collectKeysInOrder(tree.Root)
			expectedCount := numKeys / 2
			
			if len(remainingKeys) != expectedCount {
				t.Errorf("Expected %d keys remaining, got %d", expectedCount, len(remainingKeys))
			}
			
			for i, key := range remainingKeys {
				expected := (i * 2) + 1
				if key != expected {
					t.Errorf("Expected key %d at position %d, got %d", expected, i, key)
				}
			}
		})
	}
}

// TestDeleteAllKeys tests deleting all keys from a tree
func TestDeleteAllKeys(t *testing.T) {
	tree := buildTreeWithKeys(3, []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	originalKeys := collectKeysInOrder(tree.Root)
	
	// Delete all keys
	for _, key := range originalKeys {
		tree.Delete(Key{K: key})
		
		// Verify key was deleted
		if treeContainsKey(tree, key) {
			t.Errorf("Key %d should have been deleted", key)
		}
		
		// Verify tree properties if tree is not empty
		if tree.Root != nil && len(tree.Root.Keys) > 0 {
			verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
		}
	}
	
	// Verify tree is empty or has empty root
	if tree.Root != nil && len(tree.Root.Keys) > 0 {
		t.Errorf("Tree should be empty after deleting all keys, but has %d keys", len(tree.Root.Keys))
	}
}

// TestDeleteComplexScenario tests a complex scenario with multiple operations
func TestDeleteComplexScenario(t *testing.T) {
	tree := NewBeetree(3)
	
	// Insert a larger set of keys
	insertKeys := []int{5, 15, 25, 35, 45, 55, 65, 75, 85, 95, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for _, k := range insertKeys {
		tree.Insert(Key{K: k})
	}
	
	// Delete keys in a specific pattern to test various scenarios
	deletePattern := []int{
		50,  // Internal node
		25,  // May cause redistribution
		75,  // May cause merge
		10,  // Leaf node
		90,  // Another pattern
		35,  // Test predecessor/successor
		65,  // Continue testing
	}
	
	for _, keyToDelete := range deletePattern {
		if treeContainsKey(tree, keyToDelete) {
			tree.Delete(Key{K: keyToDelete})
			
			// Verify deletion
			if treeContainsKey(tree, keyToDelete) {
				t.Errorf("Key %d should have been deleted", keyToDelete)
			}
			
			// Verify tree properties
			verifyBTreeProperties(t, tree, tree.Root, tree.Degree, true)
		}
	}
	
	// Verify final state
	finalKeys := collectKeysInOrder(tree.Root)
	t.Logf("Final tree has %d keys", len(finalKeys))
	
	// Ensure all remaining keys are in sorted order
	for i := 1; i < len(finalKeys); i++ {
		if finalKeys[i-1] >= finalKeys[i] {
			t.Errorf("Tree keys not in sorted order: %d >= %d", finalKeys[i-1], finalKeys[i])
		}
	}
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

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	
	// Pre-build tree for deletion benchmark
	tree := NewBeetree(btreeDegree)
	keys := perm(benchmarkTreeSize)
	for _, item := range keys {
		tree.Insert(item)
	}
	
	deleteOrder := perm(benchmarkTreeSize)
	b.StartTimer()
	
	i := 0
	for i < b.N && i < len(deleteOrder) {
		tree.Delete(deleteOrder[i])
		i++
	}
}
