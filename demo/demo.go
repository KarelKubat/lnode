package main

import (
	"fmt"

	"github.com/KarelKubat/lnode"
)

func main() {
	var root *lnode.Node[string]

	// Initialize a linked list.
	for _, s := range []string{"the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"} {
		n := lnode.New[string](s)
		if root == nil {
			root = n
		} else {
			root.Tail().Append(n)
		}
	}

	// Visit all nodes, starting at the root.
	fmt.Println("Nodes in order:")
	root.VisitByNext(func(node *lnode.Node[string]) bool {
		fmt.Print(node.Value, " ") // Output contained value.
		return true                // Continue visiting next nodes.
	})
	fmt.Println()

	n := root.Next.Next                     // Point at the node containing string "brown"
	fmt.Println("root.Next.Next:", n.Value) // Output it
	n.Delete()                              // Remove from the chain

	n = root.Next.Next                                               // Now points at "fox" since "brown" is gone
	fmt.Println("root.Next.Next after removal of 'brown':", n.Value) // Output it

	// Visit all nodes, starting from the tail.
	fmt.Println("Nodes in reverse order:")
	root.Tail().VisitByPrev(func(node *lnode.Node[string]) bool {
		fmt.Print(node.Value, " ")
		return true
	})
	fmt.Println()

	/*
		 Output:
		 	Nodes in order:
			the quick brown fox jumps over the lazy dog
			root.Next.Next: brown
			root.Next.Next after removal of 'brown': fox
			Nodes in reverse order:
			dog lazy the over jumps fox quick the
	*/

	// Make this a circular chain
	hd := root.Head()
	tl := root.Tail()
	hd.Prev = tl
	tl.Next = hd
	fmt.Println("Is root now a circular chain:", root.Circular())

	// In a circular chain, VisitBy* will stop once every node has been seen
	hd.VisitByNext(func(node *lnode.Node[string]) bool {
		fmt.Print(node.Value, " ")
		return true
	})
	fmt.Println()

}
