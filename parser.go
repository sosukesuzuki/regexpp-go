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
func (p *Parser) ParsePattern() *Node {
	return p.parseDisjunction()
}

// https://tc39.es/ecma262/#prod-Disjunction
func (p *Parser) parseDisjunction() *Node {
	node := p.parseAlternative()
	for {
		if p.lexer.Eat(VerticalLine) {
			start := p.lexer.I
			node = &Node{
				Data: &NDisjunction{
					Left:  node,
					Right: p.parseAlternative(),
				},
				Loc: Loc{start, p.lexer.I},
			}
		} else {
			return node
		}
	}
}

/*
  https://tc39.es/ecma262/#prod-Alternative
  TODO: ループで書き直す
*/
func (p *Parser) parseAlternative() *Node {
	if p.lexer.Match(Eof) {
		return nil
	}
	start := p.lexer.I
	return &Node{
		Data: &NAlternative{
			Left: p.parseAlternative(),
			Right: p.parseTerm(),
		},
		Loc: Loc{start, p.lexer.I},
	}
}

// https://tc39.es/ecma262/#prod-Term
func (p *Parser) parseTerm() *Node {
	return p.parseAtom()
}

// https://tc39.es/ecma262/#prod-Atom
func (p *Parser) parseAtom() *Node {
	return &Node{}
}
