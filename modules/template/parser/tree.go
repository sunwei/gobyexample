package parser

type tree struct {
	root *treeNode
}

func (t *tree) String() string {
	cs := t.root.Children()
	s := ""
	for _, n := range cs {
		s += n.(Node).String()
	}
	return s
}

func newTree() *tree {
	return &tree{root: &treeNode{nodes: []TreeNode{}}}
}

type treeNode struct {
	nodes []TreeNode
}

func (n *treeNode) AppendChild(node TreeNode) {
	n.nodes = append(n.nodes, node)
}

func (n *treeNode) Children() []TreeNode {
	return n.nodes
}
