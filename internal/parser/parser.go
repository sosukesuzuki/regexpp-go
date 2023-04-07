package parser

import (
	"github.com/sosukesuzuki/regexpp-go/internal/lexer"
)

type Parser struct {
	u     bool
	lexer *lexer.Lexer
}

func NewParser(s string, u bool) Parser {
	return Parser{
		lexer: lexer.NewLexer(s, u),
	}
}

func (p *Parser) ParsePattern() {
}

func (p *Parser) parseAlternative() {
}
