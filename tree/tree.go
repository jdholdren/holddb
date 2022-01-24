package holddb

import (
	"bytes"
	"sort"
)

type node struct {
	// The prefix compressed into this node
	prefix []byte

	// Optionally store a leaf value
	leaf *leaf

	// Edges going to other nodes
	edges []*edge
}

// Searches for an edge on the node that has the correct prefix
func (n *node) getEdge(prefix []byte) (*edge, int) {
	p := prefix[0]
	l := len(n.edges)
	idx := sort.Search(l, func(i int) bool {
		return n.edges[i].label >= p
	})

	if idx >= l {
		return nil, idx
	}

	if n.edges[idx].label != p {
		return nil, idx
	}

	return n.edges[idx], idx
}

// Adds an edge to the slice in a sorted manner
func (n *node) addEdge(e *edge) {
	_, idx := n.getEdge([]byte{e.label})
	// Assume the edges are allready sorted
	n.edges = append(n.edges, nil)
	copy(n.edges[idx+1:], n.edges[idx:])
	n.edges[idx] = e
}

type edge struct {
	// Cool property of radix trees: We can just store the first by of the prefix
	// and compare that
	label byte

	// The node this edge connects to
	node *node
}

type leaf struct {
	key   string
	value interface{}
}

type Tree struct {
	root *node
}

// New creates a tree ready to accept values
func New() *Tree {
	return &Tree{
		root: &node{},
	}
}

// Insert adds a new value to the tree and returns if it was an update
func (t *Tree) Insert(key []byte, val interface{}) {
	// This holds the leaf we're going to put... somewhere
	newLeaf := &leaf{
		key:   string(key),
		value: val,
	}

	// Get the longest set of nodes that match our key
	nodes := t.longestPath(key)
	l := len(nodes)

	// If the length is one, it's an addition to the root, which is just the case
	// where we append it to the edges
	if l == 1 {
		nodes[0].addEdge(&edge{
			label: key[0],
			node: &node{
				prefix: key,
				leaf:   newLeaf,
			},
		})
		return
	}

	// We'll be mainly modifying the farthest node and its parent
	farthest := nodes[l-1]
	farthestParent := nodes[l-2]

	// Grab how much of the key was match up till the last node
	matched := []byte{}
	for i := 0; i < l-1; i++ {
		matched = append(matched, nodes[i].prefix...)
	}

	// The remaining prefix of the key to compare to the last node
	remaining := key[len(matched):len(key)]

	// We know it's a match if the remaining bytes of key are the prefix of the last node
	if bytes.Equal(remaining, farthest.prefix) {
		farthest.leaf = newLeaf
		return
	}

	// Find what we have in common
	cp := commonPrefix(remaining, farthest.prefix)
	// Find out what's different
	leafPrefix := remaining[len(cp):len(remaining)]
	oldPrefix := farthest.prefix[len(cp):len(farthest.prefix)]

	// If the old prefix is empty, then it's been entirely matched, and this is an
	// append to its edges
	if len(oldPrefix) == 0 {
		farthest.addEdge(&edge{
			label: leafPrefix[0],
			node: &node{
				prefix: leafPrefix,
				leaf:   newLeaf,
			},
		})
		return
	}

	// Something in old prefix remains, so we must split it into what it had in common
	// with the new key, and then add the new leaf node and the old "farthest" node as
	// children
	splitNode := &node{
		prefix: cp,
	}

	// We need to add our new value as a leaf to the split node
	splitNode.addEdge(&edge{
		label: leafPrefix[0],
		node: &node{
			prefix: leafPrefix,
			leaf:   newLeaf,
		},
	})

	// Add the old node to the split node
	splitNode.addEdge(&edge{
		label: oldPrefix[0],
		node: &node{
			prefix: oldPrefix,
			edges:  farthest.edges,
			leaf:   farthest.leaf,
		},
	})

	// Replace the old edge with one to the split node
	_, idx := farthestParent.getEdge(cp)
	farthestParent.edges[idx] = &edge{
		label: cp[0],
		node:  splitNode,
	}

	// panic(string(splitNode.prefix))
}

// Det
func commonPrefix(a, b []byte) []byte {
	c := []byte{}

	// We need to compare each one byte by byte to determine what the have in common.
	//
	// We can loop over either one, just check that we're not going out of bounds on the
	// other
	for i := range a {
		if i >= len(b) || a[i] != b[i] {
			break
		}

		c = append(c, a[i])
	}

	return c
}

// Returns a path of nodes that match the key the longest
//
// Should always return at least the root node.
func (t *Tree) longestPath(key []byte) []*node {
	// Start at the root
	n := t.root
	p := key

	// Remember those that have been traversed
	trav := []*node{}

	for {
		// Mark this node as traversed
		trav = append(trav, n)

		// If there's no prefix left, we're done
		if len(p) == 0 {
			break
		}

		// Get the next edge
		e, _ := n.getEdge(p)
		if e == nil {
			// Next edge couldn't be found
			break
		}

		// Set the next node to the one we found
		n = e.node

		// Lop off the prefix that we've traversed
		p = p[len(e.node.prefix):len(p)]
	}

	return trav
}
