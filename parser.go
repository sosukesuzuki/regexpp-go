package regexpp

type Parser struct {
	u     bool
	lexer *Lexer
}

func NewParser(s string, u bool) Parser {
	return Parser{
		lexer: NewLexer(s, u),
	}
}

// https://tc39.es/ecma262/#prod-Pattern
func (p *Parser) ParsePattern() Node {
	return p.ParseDisjunction()
}

// https://tc39.es/ecma262/#prod-Disjunction
func (p *Parser) ParseDisjunction() Node {
	node := p.ParseAlternative()
	for {
		if p.lexer.Eat(VerticalLine) {
			start := p.lexer.I
			node = Node{
				Data: &NDisjunction{
					Left:  node,
					Right: p.ParseAlternative(),
				},
				Loc: Loc{start, p.lexer.I},
			}
		} else {
			return node
		}
	}
}

// https://tc39.es/ecma262/#prod-Alternative
func (p *Parser) ParseAlternative() Node {
	return Node{}
}
