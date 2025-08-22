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

	anchor := lnode.New[int](0) // Node with int value zero
*/
func New[V any](value V) *Node[V] {
	return &Node[V]{
		Value: value,
	}
}

/*
Prepend adds a new node "left" of the current node. Example:

	anchor := lnode.New[int](0)
	anchor.Prepend(New[int](-1))
	// Structure:
	// -1 --- 0
	//        ^anchor
*/
func (n *Node[V]) Prepend(node *Node[V]) {
	oldPrev := n.Prev
	node.Next = n
	n.Prev = node
	node.Prev = oldPrev
	if oldPrev != nil {
		oldPrev.Next = node
	}
}

/*
Prepend adds a new node "right" of the current node. Example:

	anchor := lnode.New[int](0)
	anchor.Append(New[int](1))
	// Structure:
	// 0 --- 1
	// ^anchor
*/
func (n *Node[V]) Append(node *Node[V]) {
	oldNext := n.Next
	node.Prev = n
	n.Next = node
	node.Next = oldNext
	if oldNext != nil {
		oldNext.Prev = node
	}
}

/*
VisitByNext invokes a visitor function (callback) on the applicable node, and on all next nodes ("to the right"). The callback returns a bool indicating whether the processing should stop. When the callback returns false, no further nodes are processed.

In the case of a circular chain (see function Circular()), VisitByNext() will stop before visiting a previously seen node.

Example:

	anchor := lnode.New[int](0)
	anchor.Append(New[int](1))
	anchor.Next.Append(New[int](2))
	anchor.Next.Next.Append(New[int](3))
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
	start := n
	for n != nil {
		if !fn(n) {
			break
		}
		n = n.Next
		if n == start {
			return
		}
	}
}

/*
VisitByPrev invokes a visitor function (callback) on the applicable node, and on all previous nodes ("to the left"). The callback returns a bool indicating whether the processing should stop. When the callback returns false, no further nodes are processed.

In the case of a circular chain (see function Circular()), VisitByPrev() will will stop before visiting a previously seen node.

Example:

	anchor := lnode.New[int](0)
	anchor.Prepend(New[int](-1))
	anchor.Prev.Prepend(New[int](-2))
	anchor.Prev.Prev.Prepend(New[int](-3))
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
	start := n
	for n != nil {
		if !fn(n) {
			break
		}
		n = n.Prev
		if n == start {
			return
		}
	}
}

/*
Head returns the "leftmost" node in a chain, i.e., the node where Prev is nil. The runtime is O(N) with N being the number of nodes "to the left".

In the case of a circular chain (see function Circular()), Head() return nil.

Example:

	anchor := lnode.New[int](0)
	anchor.Prepend(New[int](-1))
	anchor.Prev.Prepend(New[int](-2))
	anchor.Prev.Prev.Prepend(New[int](-3))
	// Structure:
	// -3 --- -2 --- -1 --- 0
	//  ^head               ^anchor

	fmt.Println(anchor.Head().Value)
	// Output: -3
*/
func (n *Node[V]) Head() *Node[V] {
	start := n
	for n.Prev != nil {
		n = n.Prev
		if n == start {
			return nil
		}
	}
	return n
}

/*
Tail returns the "rightmost" node in a chain, i.e, the node where Next is nil. The runtime is O(N) with N being the number of nodes "to the right".

In the case of a circular chain (see function Circular()), Tail() will return nil.

Example:

	anchor := lnode.New[int](0)
	anchor.Append(New[int](1))
	anchor.Next.Append(New[int](2))
	anchor.Next.Next.Append(New[int](3))
	// Structure:
	// 0 --- 1 --- 2 --- 3
	// ^anchor           ^tail

	fmt.Println(anchor.Tail().Value)
	// Output: 3
*/
func (n *Node[V]) Tail() *Node[V] {
	start := n
	for n.Next != nil {
		n = n.Next
		if n == start {
			return nil
		}
	}
	return n
}

/*
Delete removes a node from the list. Example:

	anchor := lnode.New[int](0)
	anchor.Append(New[int](1))
	anchor.Next.Append(New[int](2))
	anchor.Next.Next.Append(New[int](3))
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

/*
Circular returns true when the node is part of a circular chain, else false. Example:

anchor := New[int](0)

	anchor.Append(New[int](1))
	anchor.Next.Append(New[int](2))
	anchor.Next.Next.Append(New[int](3))

	// Structure:
	// 0 --- 1 --- 2 --- 3
	// ^anchor           ^tail
	fmt.Println("is anchor in a circular chain:", anchor.Circular())  // false

	anchor.Next.Next.Next.Next = anchor
	// Structure:
	// +-----------------+
	// |                 |
	// 0 --- 1 --- 2 --- 3
	// ^anchor
	fmt.Println("is anchor in a circular chain:", anchor.Circular())  // true
*/
func (n *Node[V]) Circular() bool {
	start := n
	for n != nil {
		if n.Next == nil {
			return false
		}
		n = n.Next
		if n == start {
			return true
		}
	}
	return false
}
