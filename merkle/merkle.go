package merkle

import (
	"fmt"
	"strconv"
)

type Tree struct {
	root     *node
	leaves   []*node
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
	t.leaves = t.makeLeaves()
	t.build(t.leaves)
}

func (t *Tree) AuditProof(leafHash string) []AuditNode {
	auditTrail := make([]AuditNode, 0, 4)

	if leaf := t.findLeaf(leafHash); leaf != nil {
		if leaf.parent == nil {
			panic("Expected leaf to have a parent.")
		}
		parent := leaf.parent
		t.buildAuditTrail(&auditTrail, parent, leaf)
	}

	return auditTrail
}

func (t *Tree) buildAuditTrail(auditTrail *[]AuditNode, parent *node, child *node) {
	if parent != nil {
		if !parent.equals(child.parent) {
			panic("Parent of child is not expected parent.")
		}

		var nextChild *node
		var branch Direction

		if child.equals(parent.left) {
			nextChild = parent.right
			branch = left
		} else {
			nextChild = parent.left
			branch = right
		}

		if nextChild != nil {
			*auditTrail = append(*auditTrail, *makeAuditNode(nextChild.hashInfo, branch))
		}

		t.buildAuditTrail(auditTrail, child.parent.parent, child.parent)
	}
}

func (t *Tree) findLeaf(hash string) *node {
	for _, v := range t.leaves {
		if v.hashInfo == hash {
			return v
		}
	}
	return nil
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

func (t *Tree) VerifyAudit(auditTrail []AuditNode, leafHash string) bool {
	if len(auditTrail) == 0 {
		panic("Audit trail cannot be empty.")
	}

	testHash := leafHash

	for _, v := range auditTrail {
		switch v.branch {
		case left:
			testHash = t.hashFunc([]byte(testHash + v.hashInfo))
		case right:
			testHash = t.hashFunc([]byte(v.hashInfo + testHash))
		}
	}

	return t.root.hashInfo == testHash
}
