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
		node:    &pattern,
	}
}

func (p *Parser) ParsePattern() *regexp_ast.Pattern {
	return p.pattern
}
