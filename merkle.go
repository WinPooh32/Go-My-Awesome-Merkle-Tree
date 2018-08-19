package merkle

import (
	"fmt"
	"strconv"
)

type node struct {
	parent *node
	left   *node
	right  *node

	hashInfo string
	data     []byte
}

func (n *node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

type Tree struct {
	root     *node
	data     [][]byte
	hashFunc func([]byte) string
}

func MakeTree(hasher func([]byte) string) *Tree {
	t := Tree{hashFunc: hasher}
	return &t
}

func (t *Tree) makeLeafNode(data []byte) *node {
	n := node{hashInfo: t.hashFunc(data), data: data}
	return &n
}

func (t *Tree) makeLeaves() []*node {
	nodes := make([]*node, 0, 4)

	for _, v := range t.data {
		nodes = append(nodes, t.makeLeafNode(v))
	}

	return nodes
}

func (t *Tree) makeNode(left, right *node) *node {
	if left == nil && right == nil {
		panic("At least one node should not be a nil!")
	}

	var hash string
	n := node{}

	switch {
	case right == nil:
		hash = left.hashInfo
		left.parent = &n
	default:
		left.parent = &n
		right.parent = &n
		hash = left.hashInfo + right.hashInfo
	}

	n.hashInfo = t.hashFunc([]byte(hash))
	n.left = left
	n.right = right

	return &n
}

func (t *Tree) print(root *node, indent int, line string) {
	if root == nil {
		return
	}

	leafData := ""

	if root.isLeaf() {
		leafData = " - " + fmt.Sprintf("%v", root.data)
	}

	format := "%" + strconv.Itoa((len(root.hashInfo)+len(line))*indent+len(leafData)) + "s\n"

	t.print(root.left, indent+1, "/")
	fmt.Printf(format, line+root.hashInfo+leafData)
	t.print(root.right, indent+1, "\\")
}

func (t *Tree) Print() {
	t.print(t.root, 1, "")
}

func (t *Tree) Hash() string {
	return t.root.hashInfo
}

func (t *Tree) Insert(datas [][]byte) {
	t.data = append(t.data, datas...)
	leaves := t.makeLeaves()
	t.build(leaves)
}

func (t *Tree) build(nodes []*node) {
	if len(nodes) == 1 {
		t.root = nodes[0]
	} else {
		parents := make([]*node, 0, 4)
		length := len(nodes)

		for i := 0; i < length; i += 2 {
			var right, parent *node

			if i+1 < length {
				right = nodes[i+1]
			} else {
				right = nil
			}

			parent = t.makeNode(nodes[i], right)
			parents = append(parents, parent)
		}

		t.build(parents)
	}
}
