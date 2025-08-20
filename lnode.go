// Package lnode provides generic nodes for doubly linked lists.
package lnode

// Node is the receiver.
type Node[V any] struct {
	Value V        // Generic contained value
	Next  *Node[V] // Pointer to next node
	Prev  *Node[V] // Pointer to previous node
}

/*
New returns an initialized Node with a contained value, and with nil Next/Prev pointers. Example:

	anchor := lnode.New([int]0) // Node with int value zero
*/
func New[V any](value V) *Node[V] {
	return &Node[V]{
		Value: value,
	}
}

/*
Prepend adds a new node "left" of the current node. Example:

	anchor := lnode.New([int]0)
	anchor.Prepend(New[int]-1)
	// Structure:
	// -1 --- 0
	//        ^anchor
*/
func (n *Node[V]) Prepend(node *Node[V]) {
	node.Prev = n.Prev
	n.Prev = node
	node.Next = n
}

/*
Prepend adds a new node "right" of the current node. Example:

	anchor := lnode.New([int]0)
	anchor.Append(New[int]1)
	// Structure:
	// 0 --- 1
	// ^anchor
*/
func (n *Node[V]) Append(node *Node[V]) {
	node.Next = n.Next
	n.Next = node
	node.Prev = n
}

/*
VisitByNext invokes a visitor function (callback) on the applicable node, and on all next nodes ("to the right"). The callback returns a bool indicating whether the processing should stop. When the callback returns false, no further nodes are processed. Example:

	anchor := lnode.New([int]0)
	anchor.Append(New[int]1)
	anchor.Next.Append(New[int]2)
	anchor.Next.Next.Append(New[int]3)
	// Structure:
	// 0 --- 1 --- 2 --- 3
	// ^anchor

	anchor.VisitByNext(func (node *Node[int])bool {
	  fmt.Println(node.V)
	  return true
	})
	// output:
	// 0
	// 1
	// 2
	// 3

	anchor.VisitByNext(func (node *Node[int])bool {
	  fmt.Println(node.V)
	  return node.V < 2
	})
	// output:
	// 0
	// 1
*/
func (n *Node[V]) VisitByNext(fn func(node *Node[V]) bool) {
	for n != nil {
		if !fn(n) {
			break
		}
		n = n.Next
	}
}

/*
VisitByPrev invokes a visitor function (callback) on the applicable node, and on all previous nodes ("to the left"). The callback returns a bool indicating whether the processing should stop. When the callback returns false, no further nodes are processed. Example:

	anchor := lnode.New([int]0)
	anchor.Prepend(New[int]-1)
	anchor.Prev.Prepend(New[int]-2)
	anchor.Prev.Prev.Prepend(New[int]-3)
	// Structure:
	// -3 --- -2 --- -1 --- 0
	//                      ^anchor

	anchor.VisitByPrev(func (node *Node[int])bool {
	  fmt.Println(node.V)
	  return true
	})
	// output:
	// 0
	// -1
	// -2
	// -3
*/
func (n *Node[V]) VisitByPrev(fn func(node *Node[V]) bool) {
	for n != nil {
		if !fn(n) {
			break
		}
		n = n.Prev
	}
}

/*
Head returns the "leftmost" node in a chain, i.e., the node where Prev is nil. The runtime is O(N) with N being the number of nodes "to the left". Example:

	anchor := lnode.New([int]0)
	anchor.Prepend(New[int]-1)
	anchor.Prev.Prepend(New[int]-2)
	anchor.Prev.Prev.Prepend(New[int]-3)
	// Structure:
	// -3 --- -2 --- -1 --- 0
	//  ^head               ^anchor

	fmt.Println(anchor.Head().Value)
	// Output: -3
*/
func (n *Node[V]) Head() *Node[V] {
	for n.Prev != nil {
		n = n.Prev
	}
	return n
}

/*
Tail returns the "rightmost" node in a chain, i.e, the node where Next is nil. The runtime is O(N) with N being the number of nodes "to the right". Example:

	anchor := lnode.New([int]0)
	anchor.Append(New[int]1)
	anchor.Next.Append(New[int]2)
	anchor.Next.Next.Append(New[int]3)
	// Structure:
	// 0 --- 1 --- 2 --- 3
	// ^anchor           ^tail

	fmt.Println(anchor.Tail().Value)
	// Output: 3
*/
func (n *Node[V]) Tail() *Node[V] {
	for n.Next != nil {
		n = n.Next
	}
	return n
}

/*
Delete removes a node from the list. Example:

	anchor := lnode.New([int]0)
	anchor.Append(New[int]1)
	anchor.Next.Append(New[int]2)
	anchor.Next.Next.Append(New[int]3)
	// Structure:
	// 0 --- 1 --- 2 --- 3
	// ^anchor

	anchor.Next.Next.Delete()
	// New structure:
	// 0 --- 1 --- 3
	// ^anchor

	anchor.Tail().Delete()
	// New structure:
	// 0 --- 1
	// ^anchor
*/
func (n *Node[V]) Delete() {
	prev := n.Prev
	next := n.Next

	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
}
