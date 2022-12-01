package parser

type market struct {
	blockParsers  []Parser
	inlineParsers []Parser
	parserMap     map[Identifier]Parser
}

var m *market = nil

func defaultMarket() Market {
	if m == nil {
		m = &market{
			blockParsers:  []Parser{newHeading(), newParagraph()},
			inlineParsers: nil,
			parserMap:     map[Identifier]Parser{},
		}
		m.initParserMap()
	}

	return m
}

func (m *market) initParserMap() {
	var ps []Parser
	ps = append(ps, m.blockParsers...)
	ps = append(ps, m.inlineParsers...)

	for _, p := range ps {
		is := p.Identifiers()
		if is == nil {
			m.parserMap[""] = p
		} else {
			for _, i := range is {
				m.parserMap[i] = p
			}
		}
	}
}

func (m *market) FindParser(identifier Identifier) Parser {
	p := m.parserMap[identifier]
	if p == nil {
		p = m.defaultParser()
	}
	return p
}

func (m *market) defaultParser() Parser {
	return m.parserMap[""]
}
