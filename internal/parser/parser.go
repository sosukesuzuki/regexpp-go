package parser

import (
	"github.com/sosukesuzuki/regexpp-go/internal/lexer"
	"github.com/sosukesuzuki/regexpp-go/internal/regexp_ast"
	"github.com/sosukesuzuki/regexpp-go/internal/unicode_consts"
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

// https://tc39.es/ecma262/#prod-Pattern
func (p *Parser) ParsePattern() *regexp_ast.Node {
	return p.parseDisjunction()
}

// https://tc39.es/ecma262/#prod-Disjunction
func (p *Parser) parseDisjunction() *regexp_ast.Node {
	node := p.parseAlternative()
	for {
		if p.lexer.Eat(unicode_consts.VerticalLine) {
			start := p.lexer.I
			node = &regexp_ast.Node{
				Data: &regexp_ast.NDisjunction{
					Left:  node,
					Right: p.parseAlternative(),
				},
				Loc: regexp_ast.Loc{
					Start: start,
					End:   p.lexer.I,
				},
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
func (p *Parser) parseAlternative() *regexp_ast.Node {
	if p.lexer.Match(unicode_consts.Eof) {
		return nil
	}
	start := p.lexer.I
	return &regexp_ast.Node{
		Data: &regexp_ast.NAlternative{
			Left:  p.parseAlternative(),
			Right: p.parseTerm(),
		},
		Loc: regexp_ast.Loc{
			Start: start,
			End:   p.lexer.I,
		},
	}
}

// https://tc39.es/ecma262/#prod-Term
func (p *Parser) parseTerm() *regexp_ast.Node {
	return p.parseAtom()
}

// https://tc39.es/ecma262/#prod-Atom
func (p *Parser) parseAtom() *regexp_ast.Node {
	return &regexp_ast.Node{}
}
