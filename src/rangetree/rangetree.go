package rangetree

import (
	"log"
)

type Range struct {
	Start int
	End   int
}

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
	Start int
	End   int
}

type RangeTree struct {
	root *TreeNode
}

func NewRangeTree() *RangeTree {
	return &RangeTree{}
}

func NewTreeNode(start int, end int) *TreeNode {
	return &TreeNode{Start: start, End: end}
}

func walk(root *TreeNode, wfunc func(node *TreeNode)) {
	if root == nil {
		return
	}

	walk(root.Left, wfunc)
	wfunc(root)
	walk(root.Right, wfunc)
}

func (rt *RangeTree) Walk(wfunc func(node *TreeNode)) {
	walk(rt.root, wfunc)
}

func overlaps(n1 *TreeNode, n2 *TreeNode) bool {
	if n1 == nil || n2 == nil {
		return false
	}

	if n1.Start >= (n2.Start-1) && n1.Start <= (n2.End+1) {
		/*
		 *    |-----|	n1
		 * |-----|	n2
		 */
		return true
	}

	if n1.End >= (n2.Start-1) && n1.End <= (n2.End+1) {
		/*
		 * |-----|	n1
		 *    |-----|	n2
		 */
		return true
	}

	return false
}

func rebuild(root *TreeNode, src *TreeNode) {
	if src == nil {
		return
	}

	insert(root, src.Start, src.End) // must come first

	rebuild(root, src.Left)
	rebuild(root, src.Right)
}

func insert(root *TreeNode, start int, end int) {
	var where **TreeNode

	if start <= (root.Start - 1) {
		if end < (root.Start - 1) {
			where = &root.Left // disjoint left
		} else {
			root.Start = start // extend left
			if end > root.End {
				root.End = end // extend right
			}
		}
	} else {
		if start > (root.End + 1) { // disjoint right
			where = &root.Right
		} else {
			if end > root.End {
				root.End = end // extend right
			}
		}
	}

	if where != nil {
		if *where == nil {
			*where = NewTreeNode(start, end)
		} else {
			insert(*where, start, end)
		}
	}

	/*
	 * handle previously disjoint regjoins becoming non-disjoint due to extensions
	 * There is probably a more elegant algorithm. This is a brute force one that
	 * detaches the subtree whose root just became overlapping.  It then traverses
	 * the subtree, re-adding each range.
	 */
	if overlaps(root.Left, root) {
		left := root.Left
		root.Left = nil
		rebuild(root, left)
	}

	if overlaps(root.Right, root) {
		right := root.Right
		root.Right = nil
		rebuild(root, right)
	}
}

func (rt *RangeTree) AddRange(start int, end int) {
	root := rt.root

	if root == nil {
		rt.root = &TreeNode{Start: start, End: end}
		return
	}

	insert(root, start, end)
}

func rangewalk(root *TreeNode, start int, end int) bool {
	if root == nil {
		return false
	}

	if start >= root.Start {
		if end <= root.End {
			return true
		}
		return false
	}

	return rangewalk(root.Right, start, end)
}

/*
 * returns true if the requested range is entirely
 * contained within the tree
 */
func (rt *RangeTree) HasRange(start int, end int) bool {
	return rangewalk(rt.root, start, end)
}

func (rt *RangeTree) Dump(verbose bool) {
	if verbose {
		rt.Walk(func(node *TreeNode) {
			log.Printf("Dump: Left %v Right %v Start %v End %v\n", node.Left, node.Right, node.Start, node.End)
		})
		return
	}

	rt.Walk(func(node *TreeNode) { log.Printf("Dump: Start %v End %v\n", node.Start, node.End) })
}

func (rt *RangeTree) Check(expected []Range) bool {
	matched := true
	i := 0

	rt.Walk(func(node *TreeNode) {
		log.Printf("Checking (%d, %d)\n", node.Start, node.End)
		if i >= len(expected) {
			log.Printf("more results (%d) than expected %d\n", i, len(expected))
			matched = false
		} else if node.Start != expected[i].Start || node.End != expected[i].End {
			matched = false
			log.Printf("mismatch at result %d: expected %v\n", i, expected[i])
		}
		i++
	})

	return matched
}
