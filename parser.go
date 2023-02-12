package regexpp

type Parser struct {
	u bool
	t *Tokenizer
}

func NewParser(s string, u bool) Parser {
	return Parser{
		t: NewTokenizer(s, u),
	}
}

// https://tc39.es/ecma262/#prod-Pattern
func (p *Parser) ParsePattern() any {
	return p.ParseDisjunction()
}

// https://tc39.es/ecma262/#prod-Disjunction
func (p *Parser) ParseDisjunction() any {
	return 3
}

// https://tc39.es/ecma262/#prod-Alternative
func (p *Parser) ParseAlternative() any {
	return 3
}
