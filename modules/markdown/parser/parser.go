package parser

type Document struct {
	*tree
}

func newDocument() *Document {
	return &Document{
		newTree(&node{
			val:        &BaseBlock{s: Closed},
			firstChild: nil,
			lastChild:  nil,
			parent:     nil,
			next:       nil,
		})}
}

func Parse(lines []Line) (*Document, error) {
	d := newDocument()
	var cob Block = nil
	var parent *node

	parent = d.root
	for _, l := range lines {
	retry:
		if cob != nil {
			s := cob.Parser().Continue(l, cob)
			switch s {
			case Children:
				ob, err := openBlock(l)
				if err != nil {
					return nil, err
				}
				n := &node{val: ob}
				parent = parent.lastChild
				parent.AppendChild(n)
				cob = ob
				continue
			case Close:
				cob.Close()
				pb := parent.val.(Block)
				if pb.IsOpen() {
					cob = pb
					parent = parent.parent
				} else {
					cob = nil
				}

				goto retry
			case Continue:
				continue
			}
		}
		ob, err := openBlock(l)
		if err != nil {
			return nil, err
		}
		cob = ob
		parent.AppendChild(&node{val: cob})
	}
	cob.Close()

	return d, nil
}

func openBlock(l Line) (Block, error) {
	market := defaultMarket()
	p := market.FindParser(Identifier(l.StartChar()))
	return p.Open(l)
}
