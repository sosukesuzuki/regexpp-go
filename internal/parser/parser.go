package parser

import (
	"github.com/sosukesuzuki/regexpp-go/internal/lexer"
	"github.com/sosukesuzuki/regexpp-go/internal/regexp_ast"
)

type Parser struct {
	u       bool
	lexer   *lexer.Lexer
	pattern *regexp_ast.Pattern
	node    regexp_ast.Node
}

func NewParser(s string, u bool) Parser {
	pattern := regexp_ast.Pattern{}
	return Parser{
		lexer:   lexer.NewLexer(s, u),
		pattern: &pattern,
		node:    nil,
	}
}

func (p *Parser) ParsePattern() *regexp_ast.Pattern {
	p.consumePattern()
	return p.pattern
}

func (p *Parser) consumePattern() {
	start := p.lexer.I
	p.onPatternEnter(start)
	p.onPatternLeave(start, p.lexer.I)
}

func (p *Parser) onPatternEnter(start int) {
	p.node = &regexp_ast.Pattern{
		Alternatives: []regexp_ast.Alternative{},
		Loc: regexp_ast.Loc{
			Start: start,
			End:   -1,
		},
	}
}

func (p *Parser) onPatternLeave(start int, end int) {
	p.node.SetEnd(end)
}
