package merkle

import "fmt"

type Direction int

const (
	left Direction = iota
	right
	oldRoot
)

type AuditNode struct {
	branch   Direction
	hashInfo string
}

func makeAuditNode(hash string, branch Direction) *AuditNode {
	if branch < left || branch > oldRoot {
		panic(fmt.Sprintf("Value %d is out of Direction type range", branch))
	}
	a := AuditNode{hashInfo: hash, branch: branch}
	return &a
}
