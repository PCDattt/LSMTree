package main

import (
	"strings"
	"strconv"
)

type Node struct {
	Key int
	Value string
	LeftNode *Node
	RightNode *Node
}

func NodeMarshal(n Node) string {
	return strconv.Itoa(n.Key) + "," + n.Value
}

func NodeUnmarshal(s string) (Node,error) {
	parts := strings.Split(s, ",")
	if len(parts) == 2 {
		k, err := strconv.Atoi(parts[0])
		if err != nil {
			return Node{}, err
		}
		v := strings.Trim(parts[1], "\n")
		return Node{
			Key: k,
			Value: v,
		}, nil
	}
	return Node{}, nil
}

func NewNode(k int, v string) *Node {
	return &Node{
		Key: k,
		Value: v,
		LeftNode: nil,
		RightNode: nil,
	}
}

type BinaryTree struct {
	Head *Node
}

func NewBinaryTree() BinaryTree {
	return BinaryTree{}
}

func (bt *BinaryTree) Clear() {
	bt.Head = nil
}

func (bt *BinaryTree) AddNode(n *Node) {
	if bt.Head == nil {
		bt.Head = n
		return
	}
	tempNode := bt.Head
	for {
		if n.Key > tempNode.Key {
			if tempNode.RightNode == nil {
				tempNode.RightNode = n
				return
			} else {
				tempNode = tempNode.RightNode
			}
		} else {
			if tempNode.LeftNode == nil {
				tempNode.LeftNode = n
				return
			} else {
				tempNode = tempNode.LeftNode
			}
		}
	}
}

func (bt *BinaryTree) FindNode(n *Node, k int) *Node {
	if n == nil {
		return nil
	}
	if k == n.Key {
		return n
	}
	if k < n.Key {
		n := bt.FindNode(n.LeftNode, k)
		if n != nil {
			return n
		}
	} else {
		n := bt.FindNode(n.RightNode, k)
		if n != nil {
			return n
		}
	}
	return nil
}

func (bt *BinaryTree) InorderTraversal(n *Node, arr *[]Node) {
	if n == nil {
		return
	}
	if n.LeftNode != nil {
		bt.InorderTraversal(n.LeftNode, arr)
	}
	*arr = append(*arr, *n)
	if n.RightNode != nil {
		bt.InorderTraversal(n.RightNode, arr)
	}
}
