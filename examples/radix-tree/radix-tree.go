package main

import (
	"fmt"
	"sort"
)

func main() {
	t := New()
	type kv struct {
		k string
		v string
	}
	kvs := []kv{
		{"/content/a.md", "aaa"},
		{"/content/b.md", "bbb"},
		{"/content/c.md", "ccc"},
	}
	for _, kv := range kvs {
		t.Insert(kv.k, kv.v)
	}
	t.Print()
}

type edge struct {
	name  string
	start *node
	end   *node
}

type edges []*edge

type node struct {
	val      any
	prefix   *edge
	suffixes edges
}

type Tree struct {
	root *node
}

// New returns an empty Tree
func New() *Tree {
	return &Tree{root: &node{
		val:      "root",
		prefix:   nil,
		suffixes: edges{},
	}}
}

// Insert is used to add/update a value to radix tree
func (t *Tree) Insert(k string, v interface{}) bool {
	var parent *node
	n := t.root
	path := k

	if len(path) == 0 {
		panic("empty key not supported yet.")
	}

	for {
		// Look for the edge
		parent = n
		e := n.getEdge(path)

		// No same prefix edge found, create one
		if e == nil {
			e := newNodeEdge()
			e.name = path
			e.start = parent
			e.end.val = v

			parent.addEdge(e)
			return true
		}

		// Determine the longest prefix of the search key on match
		commonPrefix := longestPrefix(path, e.name)
		// Edge found with fully overlap
		// Look into the right depth
		if commonPrefix == len(e.name) {
			path = path[commonPrefix:]
			n = e.end
			continue
		}

		// Right depth found
		// Split the current node
		// Create common edge
		commonEdge := newNodeEdge()
		commonEdge.name = path[:commonPrefix]
		commonEdge.start = e.start

		e.start.addEdge(commonEdge)
		e.start.delEdge(e)
		commonEdge.end.addEdge(e)

		// Update edge name with uncommon part
		e.name = e.name[commonPrefix:]

		// Create the new joined one
		freshEdge := newNodeEdge()
		freshEdge.name = path[commonPrefix:]
		freshEdge.start = commonEdge.end
		freshEdge.end.val = v

		commonEdge.end.addEdge(freshEdge)
		return true
	}

	return false
}

func newNodeEdge() *edge {
	e := &edge{
		name:  "",
		start: nil,
		end: &node{
			val:      nil,
			prefix:   nil,
			suffixes: edges{},
		},
	}
	e.end.prefix = e
	return e
}

func (n *node) getEdge(path string) *edge {
	return findEdgeWithSamePrefix(
		getFirstByte(path), n.suffixes)
}

func (n *node) addEdge(e *edge) {
	e.start = n
	n.suffixes = append(n.suffixes, e)
	sortAscending(n.suffixes)
}

func (n *node) delEdge(e *edge) {
	if e.start == n {
		var pos int
		for i, item := range n.suffixes {
			if item == e {
				pos = i
				break
			}
		}
		n.suffixes = append(n.suffixes[:pos],
			n.suffixes[pos+1:]...)
	}
}

func sortAscending(es edges) {
	sort.Slice(es, func(i, j int) bool {
		return getFirstByte(
			es[i].name) < getFirstByte(es[j].name)
	})
}

func getFirstByte(v string) byte {
	return v[0]
}

func findEdgeWithSamePrefix(
	firstByte byte, es edges) *edge {
	for _, e := range es {
		if getFirstByte(e.name) == firstByte {
			return e
		}
	}
	return nil
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	max := len(k1)
	if l := len(k2); l < max {
		max = l
	}
	var i int
	for i = 0; i < max; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

func (t *Tree) Print() {
	recursivePrint(t.root)
}

func recursivePrint(n *node) {
	printNode(n)
	for _, e := range n.suffixes {
		recursivePrint(e.end)
	}
}

func printNode(n *node) {
	prefix := ""
	if n.prefix != nil {
		prefix = n.prefix.name
	}
	fmt.Printf("node: prefix %s, value %s\n",
		prefix, n.val)
}
