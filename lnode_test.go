package lnode

import "testing"

func TestAppend(t *testing.T) {
	start := New[int](0)
	n := start
	for i := 1; i < 5; i++ {
		n.Append(New[int](i))
		n = n.Next
	}

	expect := 0
	start.VisitByNext(func(node *Node[int]) bool {
		if node.Value != expect {
			t.Errorf("VisitByNext: got Value %d, want %d", node.Value, expect)
		}
		expect++
		return true
	})

	n = n.Head()
	for i := range 5 {
		if n.Value != i {
			t.Errorf("From Head: got Value %d, want %d", n.Value, i)
		}
		n = n.Next
	}
}

func TestPrepend(t *testing.T) {
	end := New[int](0)
	n := end
	for i := 1; i < 5; i++ {
		n.Prepend(New[int](i))
		n = n.Prev
	}

	expect := 0
	end.VisitByPrev(func(node *Node[int]) bool {
		if node.Value != expect {
			t.Errorf("VisitByPrev: got Value %d, want %d", node.Value, expect)
		}
		expect++
		return true
	})

	n = n.Tail()
	for i := range 5 {
		if n.Value != i {
			t.Errorf("From Tail: got Value %d, want %d", n.Value, i)
		}
		n = n.Prev
	}
}

func TestDelete(t *testing.T) {
	var anchor *Node[int]

	mkList := func() {
		anchor = nil
		var n *Node[int]
		for i := range 10 {
			n = New[int](i)
			if i == 0 {
				anchor = n
			} else {
				anchor.Tail().Append(n)
			}
		}
	}
	expect := func(desc string, nrs []int) {
		n := anchor
		for _, nr := range nrs {
			if n.Value != nr {
				t.Errorf("TestDelete: %s: unmatched number in %v, got %d, want %d", desc, nrs, n.Value, nr)
			}
			n = n.Next
		}

	}

	// Single delete in the middle
	mkList()
	anchor.Next.Delete()
	expect("single middle delete", []int{0, 2, 3, 4, 5, 6, 7, 8, 9})

	// Multiple deletes
	mkList()
	anchor.Next.Next.Next.Delete()
	anchor.Next.Next.Next.Delete()
	anchor.Next.Next.Next.Delete()
	expect("triple middle delete", []int{0, 1, 2, 6, 7, 8, 9})

	// Delete start
	mkList()
	nextAnchor := anchor.Next
	anchor.Delete()
	anchor = nextAnchor
	expect("anchor delete", []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	if anchor.Prev != nil {
		t.Errorf("TestDelete: anchor.Prev in %v is not nil", anchor)
	}

	// Delete end
	mkList()
	anchor.Next.Next.Next.Next.Next.Next.Next.Next.Next.Delete()
	expect("tail delete", []int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	if anchor.Tail().Next != nil {
		t.Errorf("TestDelete: anchor.Tail().Next in %v is not nil", anchor.Tail())
	}
}

func TestVisitedAndHeadTail(t *testing.T) {
	anchor := New[int](0)
	anchor.Append(New[int](1))
	anchor.Next.Append(New[int](2))
	anchor.Next.Next.Append(New[int](3))
	if got := anchor.Circular(); got != false {
		t.Errorf("Before closing the loop: anchor.Circular() returned true")
	}
	if hd := anchor.Head(); hd != anchor {
		t.Errorf("Before closing the loop: anchor.Head() = %v, want %v", hd, anchor)
	}
	if tl := anchor.Tail(); tl != anchor.Next.Next.Next {
		t.Errorf("Before closing the loop: anchor.Head() = %v, want %v", tl, anchor.Next.Next.Next)
	}

	anchor.Next.Next.Next.Next = anchor
	anchor.Prev = anchor.Next.Next.Next
	if got := anchor.Circular(); got != true {
		t.Errorf("After closing the loop: anchor.Circular() returned false")
	}
	if hd := anchor.Head(); hd != nil {
		t.Errorf("Before closing the loop: anchor.Head() = %v, want nil", hd)
	}
	if tl := anchor.Tail(); tl != nil {
		t.Errorf("Before closing the loop: anchor.Head() = %v, want nil", tl)
	}
}
