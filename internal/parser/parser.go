package parser

import (
	"errors"

	"github.com/sosukesuzuki/regexpp-go/internal/lexer"
	"github.com/sosukesuzuki/regexpp-go/internal/regexp_ast"
	"github.com/sosukesuzuki/regexpp-go/internal/unicode_consts"
)

type Parser struct {
	u       bool
	lexer   *lexer.Lexer
	pattern *regexp_ast.Pattern
	node    regexp_ast.Node
	errors  []error
}

func NewParser(s string, u bool) Parser {
	return Parser{
		lexer:   lexer.NewLexer(s, u),
		pattern: nil,
		node:    nil,
	}
}

func (p *Parser) ParsePattern() (*regexp_ast.Pattern, error) {
	p.consumePattern()
	return p.pattern, errors.Join(p.errors...)
}

func (p *Parser) raise(msg string) {
	p.errors = append(p.errors, &ParserError{
		msg: msg,
		err: nil,
	})
}

//------------------------------------------------------------------------------
// Pattern
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Pattern
//------------------------------------------------------------------------------

func (p *Parser) consumePattern() {
	start := p.lexer.I
	p.onPatternEnter(start)
	p.consumeDisjunction()
	p.onPatternLeave(start, p.lexer.I)
}

func (p *Parser) onPatternEnter(start int) {
	pattern := &regexp_ast.Pattern{
		Alternatives: []*regexp_ast.Alternative{},
		Loc: regexp_ast.Loc{
			Start: start,
			End:   -1,
		},
	}
	p.node = pattern
	p.pattern = pattern
}

func (p *Parser) onPatternLeave(start int, end int) {
	p.node.SetEnd(end)
}

//------------------------------------------------------------------------------
// Disjunction
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Disjunction
//------------------------------------------------------------------------------

func (p *Parser) consumeDisjunction() {
	start := p.lexer.I
	p.onDisjunctionEnter(start)

	i := 0
	for {
		p.consumeAlternative(i)
		i = i + 1
		if !p.lexer.Eat(unicode_consts.VerticalLine) {
			break
		}
	}

	p.onDisjunctionLeave(start, p.lexer.I)
}

func (p *Parser) onDisjunctionEnter(start int) {
	// do nothing
}

func (p *Parser) onDisjunctionLeave(start int, end int) {
	// do nothing
}

//------------------------------------------------------------------------------
// Alternative
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Alternative
//------------------------------------------------------------------------------

func (p *Parser) consumeAlternative(index int) {
	start := p.lexer.I

	p.onAlternativeEnter(start)

	for {
		if p.lexer.CP == -1 || !p.consumeTerm() {
			break
		}
	}

	p.onAlternativeLeave(start, p.lexer.I)
}

func (p *Parser) onAlternativeEnter(start int) {
	parent, ok := p.node.(*regexp_ast.Pattern)
	if !ok {
		p.raise("The parent of Alternative must be Pattern")
	}
	alt := &regexp_ast.Alternative{
		Elements: []regexp_ast.Element{},
		Parent:   parent,
		Loc: regexp_ast.Loc{
			Start: start,
			End:   -1,
		},
	}
	p.node = alt
	parent.Alternatives = append(parent.Alternatives, alt)
}

func (p *Parser) onAlternativeLeave(start int, end int) {
	p.node.SetEnd(end)
	p.node = p.node.GetParent()
}

//------------------------------------------------------------------------------
// Term
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Term
//------------------------------------------------------------------------------

func (p *Parser) consumeTerm() bool {
	return p.consumeAssertion() || (p.consumeAtom() && p.consumeOptionalQuantifier())
}

//------------------------------------------------------------------------------
// Assertion
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Assertion
//------------------------------------------------------------------------------

func (p *Parser) consumeAssertion() bool {
	return false
}

//------------------------------------------------------------------------------
// Atom
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
//------------------------------------------------------------------------------

func (p *Parser) consumeAtom() bool {
	return p.consumePatternCharacter() ||
		p.consumeDot() ||
		p.consumeReverseSolidusAtomEscape() ||
		p.consumeCharacterClass() ||
		p.consumeUncapturingGroup() ||
		p.consumeCapturingGroup()
}

//------------------------------------------------------------------------------
// Quantifier
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Quantifier
//------------------------------------------------------------------------------

func (p *Parser) consumeOptionalQuantifier() bool {
	p.consumeQuantifier()
	return true
}

func (p *Parser) consumeQuantifier() bool {
	return false
}

//------------------------------------------------------------------------------
// PatternCharacter
// https://tc39.es/ecma262/multipage/text-processing.html#prod-PatternCharacter
//------------------------------------------------------------------------------

func (p *Parser) consumePatternCharacter() bool {
	return false
}

//------------------------------------------------------------------------------
// .
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
//------------------------------------------------------------------------------

func (p *Parser) consumeDot() bool {
	return false
}

// ------------------------------------------------------------------------------
// \ AtomEscape
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
// ------------------------------------------------------------------------------
func (p *Parser) consumeReverseSolidusAtomEscape() bool {
	return false
}

// ------------------------------------------------------------------------------
// CharacterClass
// https://tc39.es/ecma262/multipage/text-processing.html#prod-CharacterClass
// ------------------------------------------------------------------------------
func (p *Parser) consumeCharacterClass() bool {
	return false
}

// ------------------------------------------------------------------------------
// ( GroupSpecifier Disjunction )
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
// ------------------------------------------------------------------------------
func (p *Parser) consumeUncapturingGroup() bool {
	return false
}

// ------------------------------------------------------------------------------
// (?: Disjunction )
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
// ------------------------------------------------------------------------------
func (p *Parser) consumeCapturingGroup() bool {
	return false
}
