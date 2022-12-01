package parser

type node struct {
	val        any
	firstChild *node
	lastChild  *node
	parent     *node
	next       *node
}

func (n *node) AppendChild(child *node) {
	if n.firstChild == nil {
		n.firstChild = child
		child.next = nil
	} else {
		last := n.lastChild
		last.next = child
	}
	child.parent = n
	n.lastChild = child
}

type tree struct {
	root *node
}

func newTree(r *node) *tree {
	return &tree{root: r}
}

func (t *tree) Walk(walker Walker) {
	walkNode(t.root, walker)
}

func walkNode(n *node, walker Walker) WalkStatus {
	status := walker(n.val, WalkIn)
	if status != WalkStop {
		for c := n.firstChild; c != nil; c = c.next {
			if s := walkNode(c, walker); s == WalkStop {
				return WalkStop
			}
		}
	}
	return walker(n.val, WalkOut)
}
