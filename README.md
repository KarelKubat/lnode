# lnode

Package `lnode` provides generics for nodes in a doubly linked list.

<!-- toc -->
- [Synopsis](#synopsis)
- [Description](#description)
<!-- /toc -->

## Synopsis

```go
package main

import (
    "fmt"

    "github.com/KarelKubat/lnode"
)

func main() {
    var root *lnode.Node[string]

    // Initialize a linked list.
    for _, s := range []string{
    		"the", "quick", "brown", "fox",
     		"jumps", "over", "the", "lazy", "dog",
    } {
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

    n := root.Next.Next  // Point at the node containing string "brown"
    n.Delete()           // Remove from the chain

    n = root.Next.Next   // Repoint, now points at "fox"
    fmt.Println("root.Next.Next after removal of 'brown':", n.Value)

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
}
```

## Description

Package `lnode` provides the definition of nodes in a doubly linked list, with a few helper functions. The base data type is a struct with the fields `Value` (a generic) and with `Prev` and `Next`, pointers to the "left" or "right".

The list is extended using `Append()` or `Prepend`(), which insert a node resp. before or after the node in question. The `Next` and `Prev` pointers are adjusted accordingly. Similarly, nodes can be removed from the list using `Delete()` which also adjusts `Prev` and `Next`.

Finally there are some trivial helpers:

- `VisitByNext()` or `VisitByPrev()` "walk" the list and invoke a callback.
- `Head()` and `Tail()` return the first, cq. last node in a chain. These iterate from the indicated node, so that they run on O(N) time.
- `Circular()` returns `true` when nodes are arranged in a circular chain (in which case, `Head()` and `Tail()` will return `nil`). This function runs in O(N) time.

*Package `lnode` is not thread-safe. The caller must ensure that concurrent updates to nodes (e.g., `lnode.Append()` calls hitting the same node) are mutex-protected.*
