package merkle

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func assertPanic(t *testing.T, testCase string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic in case: %s", testCase)
		}
	}()
	f()
}

func md5Wrap(data []byte) string {
	sum := md5.Sum(data)
	return fmt.Sprintf("%X", sum)
}

func makeTestTree() *Tree {
	tree := MakeTree(md5Wrap)
	return tree
}

func makeTestLeaf(testData []byte) (*Tree, *node) {
	tree := makeTestTree()
	n := tree.makeLeafNode(testData)
	return tree, n
}

func TestMakeTree(t *testing.T) {
	tree := MakeTree(md5Wrap)

	if tree == nil {
		t.Fatal()
	}

	if tree.data != nil {
		t.Fatal()
	}

	if tree.leaves != nil {
		t.Fatal()
	}

	if tree.root != nil {
		t.Fatal()
	}

	if tree.hashFunc([]byte("test")) != md5Wrap([]byte("test")) {
		t.Fatal()
	}
}

func testArrEq(a, b []byte) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestMakeLeafNode(t *testing.T) {

	cases := [][]byte{
		[]byte{},
		[]byte{1},
		[]byte{1, 2},
		[]byte{1, 2, 3},
	}

	for _, v := range cases {
		tree, n := makeTestLeaf(v)

		if tree == nil {
			t.Fatal()
		}

		if n == nil {
			t.Fatal()
		}

		if n.data == nil {
			t.Fatal()
		}

		if !testArrEq(n.data, v) {
			t.Fatal()
		}

		if n.hashInfo == "" {
			t.Fatal()
		}

		if tree.hashFunc(v) != n.hashInfo {
			t.Fatal()
		}
	}

	assertPanic(t, "nil data", func() { makeTestLeaf(nil) })
}

func TestIsLeaf(t *testing.T) {
	tree := makeTestTree()
	n := tree.makeLeafNode([]byte{})

	if !n.isLeaf() {
		t.Fatal()
	}
}

func leafEqualsCases(t *testing.T) {
	tree := makeTestTree()
	n1 := tree.makeLeafNode([]byte{1, 2, 3})
	n2 := tree.makeLeafNode([]byte{})
	n3 := tree.makeLeafNode([]byte{1, 2, 3})

	if !n1.equals(n1) {
		t.Fatal()
	}

	if n1.equals(n2) {
		t.Fatal()
	}

	if !n1.equals(n3) || !n3.equals(n1) {
		t.Fatal()
	}
}

func TestEquals(t *testing.T) {
	leafEqualsCases(t)
}

func TestMakeLeaves(t *testing.T) {
	testLeaves := func(leaves []*node, dataSet [][]byte) {
		for i, v := range leaves {
			if !v.isLeaf() {
				t.Fatal()
			}

			if !testArrEq(v.data, dataSet[i]) {
				t.Fatal()
			}
		}
	}

	func() {
		tree := makeTestTree()
		leaves := tree.makeLeaves()
		testLeaves(leaves, [][]byte{})
	}()

	cases := []([][]byte){
		{},
		{{1}},
		{{1}, {2}},
		{{1}, {2}, {3}},
	}

	for _, v := range cases {
		tree := makeTestTree()
		tree.data = v
		leaves := tree.makeLeaves()
		testLeaves(leaves, v)
	}
}
