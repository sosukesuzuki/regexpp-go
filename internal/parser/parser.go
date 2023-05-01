package parser

import (
	"errors"
	"math"

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
	state   *struct {
		lastIntValue int
		lastMaxValue int
		lastMinValue int
	}
}

func NewParser(s string, u bool) Parser {
	return Parser{
		lexer:   lexer.NewLexer(s, u),
		pattern: nil,
		node:    nil,
		state: &struct {
			lastIntValue int
			lastMaxValue int
			lastMinValue int
		}{
			lastIntValue: 0,
			lastMaxValue: 0,
			lastMinValue: 0,
		},
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

// Quantifier::
//
//	QuantifierPrefix
//	QuantifierPrefix `?`
//
// QuantifierPrefix::
//
//	`*`
//	`+`
//	`?`
//	`{` DecimalDigits `}`
//	`{` DecimalDigits `,}`
//	`{` DecimalDigits `,` DecimalDigits `}`
func (p *Parser) consumeQuantifier() bool {
	start := p.lexer.I
	min := 0
	max := 0
	greety := false

	if p.lexer.Eat(unicode_consts.Asterisk) {
		min = 0
		max = math.MaxInt
	} else if p.lexer.Eat(unicode_consts.PlusSign) {
		min = 1
		max = math.MaxInt
	} else if p.lexer.Eat(unicode_consts.QuestionMark) {
		min = 0
		max = 1
	} else if p.eatBracedQuantifier() {
		min = p.state.lastMinValue
		max = p.state.lastMaxValue
	} else {
		return false
	}

	greety = !p.lexer.Eat(unicode_consts.QuestionMark)

	p.onQuantifier(start, p.lexer.I, min, max, greety)

	return false
}

func (p *Parser) onQuantifier(start int, end int, min int, max int, greety bool) bool {
	switch parent := p.node.(type) {
	case *regexp_ast.Alternative:
		{
			// Replace the last element (pop)
			element := parent.Elements[len(parent.Elements)-1]
			elements := parent.Elements[:len(parent.Elements)-1]
			parent.Elements = elements

			if quantifiable, ok := element.(regexp_ast.QuantifiableElement); ok {
				q := &regexp_ast.Quantifier{
					Parent: parent,
					Loc: regexp_ast.Loc{
						Start: start,
						End:   end,
					},
					Greety:  greety,
					Min:     min,
					Max:     max,
					Element: quantifiable,
				}
				parent.Elements = append(parent.Elements, q)
				if node, ok := quantifiable.(regexp_ast.Node); ok {
					node.SetParent(q)
				} else {
					p.raise("Maybe unureaced")
					return false
				}
				return true
			}
			return false
		}
	default:
		p.raise("The parent of Quantifier must  be Alternative")
		return false
	}
}

func (p *Parser) eatBracedQuantifier() bool {
	start := p.lexer.I
	if p.lexer.Eat(unicode_consts.LeftCurlyBracket) {
		p.state.lastMinValue = 0
		p.state.lastMaxValue = math.MaxInt
		if digit := p.eatDecimalDigits(); digit != -1 {
			p.state.lastMinValue = digit
			p.state.lastMaxValue = digit
			if p.lexer.Eat(unicode_consts.Comma) {
				if secondDigit := p.eatDecimalDigits(); secondDigit != -1 {
					p.state.lastMaxValue = secondDigit
				}
			}
			if p.lexer.Eat(unicode_consts.RightCurlyBracket) {
				return true
			}
			p.raise("Imcomplete Quantifier")
			p.lexer.Rewind(start)
		}
	}
	return false
}

//------------------------------------------------------------------------------
// PatternCharacter
// https://tc39.es/ecma262/multipage/text-processing.html#prod-PatternCharacter
//------------------------------------------------------------------------------

func (p *Parser) consumePatternCharacter() bool {
	start := p.lexer.I
	cp := p.lexer.CP
	if cp != -1 && !isSyntaxCharacter(cp) {
		p.lexer.Next()
		p.onCharacter(start, p.lexer.I, cp)
		return true
	}
	return false
}

//------------------------------------------------------------------------------
// .
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
//------------------------------------------------------------------------------

func (p *Parser) consumeDot() bool {
	if p.lexer.Eat(unicode_consts.FullStop) {
		p.onAnyCharacterSet(p.lexer.I-1, p.lexer.I)
		return true
	}
	return false
}

func (p *Parser) onAnyCharacterSet(start int, end int) {
	switch parent := p.node.(type) {
	case *regexp_ast.Alternative:
		parent.Elements = append(parent.Elements, &regexp_ast.AnyCharacterSet{
			Parent: parent,
			Loc: regexp_ast.Loc{
				Start: start,
				End:   end,
			},
		})
	default:
		p.raise("The parent of AnyCharacterSet must be Alternative")
	}
}

// ------------------------------------------------------------------------------
// \ AtomEscape
// https://tc39.es/ecma262/multipage/text-processing.html#prod-Atom
// ------------------------------------------------------------------------------
func (p *Parser) consumeReverseSolidusAtomEscape() bool {
	return false
}

// ------------------------------------------------------------------------------
// CharacterClass ::
//
//	[ [lookahead != ^] ClassRanges ]
//	[ ^ ClassRanges ]
//
// https://tc39.es/ecma262/multipage/text-processing.html#prod-CharacterClass
// ------------------------------------------------------------------------------
func (p *Parser) consumeCharacterClass() bool {
	start := p.lexer.I
	if p.lexer.Eat(unicode_consts.LeftSquareBracket) {
		negate := p.lexer.Eat(unicode_consts.CircumflexAccent)
		p.onCharacterClassEnter(start, negate)
		p.consumeClassRanges()
		if !p.lexer.Eat(unicode_consts.RightSquareBracket) {
			p.raise("Unterminated character class")
		}
		p.onCharacterClassLeave(start, p.lexer.I, negate)
		return true
	}
	return false
}

func (p *Parser) onCharacterClassEnter(start int, negate bool) {}

func (p *Parser) onCharacterClassLeave(start int, end int, negate bool) {}

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

// ------------------------------------------------------------------------------
// SourceCharacter
// https://tc39.es/ecma262/multipage/ecmascript-language-source-code.html#prod-SourceCharacter
// ------------------------------------------------------------------------------
func (p *Parser) onCharacter(start int, end int, value int) {
	switch parent := p.node.(type) {
	case *regexp_ast.Alternative:
		parent.Elements = append(parent.Elements, &regexp_ast.Character{
			Parent: parent,
			Value:  value,
			Loc: regexp_ast.Loc{
				Start: start,
				End:   end,
			},
		})
	case *regexp_ast.CharacterClass:
		parent.Elements = append(parent.Elements, &regexp_ast.Character{
			Parent: parent,
			Value:  value,
			Loc: regexp_ast.Loc{
				Start: start,
				End:   end,
			},
		})
	default:
		p.raise("The parent of Character must be Alternative or CharacterClass")
	}
}

// ------------------------------------------------------------------------------
// SyntaxCharacter
// https://tc39.es/ecma262/multipage/text-processing.html#prod-SyntaxCharacter
// ------------------------------------------------------------------------------
func isSyntaxCharacter(cp int) bool {
	return cp == unicode_consts.CircumflexAccent ||
		cp == unicode_consts.DollarSign ||
		cp == unicode_consts.ReverseSolidus ||
		cp == unicode_consts.FullStop ||
		cp == unicode_consts.Asterisk ||
		cp == unicode_consts.PlusSign ||
		cp == unicode_consts.QuestionMark ||
		cp == unicode_consts.LeftParenthesis ||
		cp == unicode_consts.RightParenthesis ||
		cp == unicode_consts.LeftSquareBracket ||
		cp == unicode_consts.RightSquareBracket ||
		cp == unicode_consts.LeftCurlyBracket ||
		cp == unicode_consts.RightCurlyBracket ||
		cp == unicode_consts.VerticalLine
}

// ------------------------------------------------------------------------------
// DecimalDigit, DecimalDigits
// https://tc39.es/ecma262/multipage/notational-conventions.html#prod-grammar-notation-DecimalDigit
//
// DecimalDigit :: one of
//   0 1 2 3 4 5 6 7 8 9
//
// DegimalDigits ::
//   DecimalDigit
//   DecimalDigits DecimalDigit
//
// ------------------------------------------------------------------------------

// Eat DecimalDigits. Returns int value that is eaten last time. If eating is failed, returns -1.
func (p *Parser) eatDecimalDigits() int {
	start := p.lexer.I
	lastInt := 0
	for {
		if !unicode_consts.IsDecimalDigit(p.lexer.CP) {
			break
		}
		lastInt = 10*lastInt + unicode_consts.DecimalToDigit(p.lexer.CP)
		p.lexer.Next()
	}
	if p.lexer.I != start {
		return lastInt
	} else {
		return -1
	}
}

// ------------------------------------------------------------------------------
// ClassRanges ::
//
//	[empty]
//	NoemptyClassRanges
//
// https://tc39.es/ecma262/multipage/text-processing.html#prod-ClassRanges
// ------------------------------------------------------------------------------
func (p *Parser) consumeClassRanges() {
	for {
		rangeStart := p.lexer.I
		if !p.consumeClassAtom() {
			break
		}
		min := p.state.lastIntValue

		if !p.lexer.Eat(unicode_consts.HyphenMinus) {
			continue
		}

		p.onCharacter(p.lexer.I-1, p.lexer.I, unicode_consts.HyphenMinus)

		if !p.consumeClassAtom() {
			break
		}
		max := p.state.lastIntValue

		p.onCharacterClassRange(rangeStart, p.lexer.I, min, max)
	}
}

func (p *Parser) onCharacterClassRange(start int, end int, min int, max int) {
	switch parent := p.node.(type) {
	case *regexp_ast.CharacterClass:
		three := parent.Elements[len(parent.Elements)-3 : len(parent.Elements)]
		if len(three) != 3 {
			p.raise("UnknownError")
		}

		minEl := three[0]
		hyphenEl := three[1]
		maxEl := three[2]
		parent.Elements = parent.Elements[:len(parent.Elements)-3]
		if minEl == nil || hyphenEl == nil || maxEl == nil {
			p.raise("UnknwonError")
		}

		if minChar, ok := minEl.(*regexp_ast.Character); ok {
			if hyphenChar, ok := hyphenEl.(*regexp_ast.Character); ok && hyphenChar.Value == unicode_consts.HyphenMinus {
				if maxChar, ok := maxEl.(*regexp_ast.Character); ok {
					node := &regexp_ast.CharacterClassRange{
						Parent: parent,
						Loc: regexp_ast.Loc{
							Start: start,
							End:   end,
						},
						Min: minChar,
						Max: maxChar,
					}
					parent.Elements = append(parent.Elements, node)
					return
				}
			}
		}
		p.raise("UnkownError")
	default:
		p.raise("The parentof CharacterClassRange must be CharacterClass")
	}
}

// ------------------------------------------------------------------------------
// ClassAtom::
//
//	-
//	ClassAtomNoDash
//
// ClassAtomNoDash
//
//	SourceCharacter but not one of / or ] or -
//	\ ClassEscape
//
// https://tc39.es/ecma262/multipage/text-processing.html#prod-ClassAtom
// ------------------------------------------------------------------------------
func (p *Parser) consumeClassAtom() bool {
	start := p.lexer.I
	cp := p.lexer.CP

	if cp != -1 && cp != unicode_consts.ReverseSolidus && cp != unicode_consts.RightSquareBracket {
		p.lexer.Next()
		p.state.lastIntValue = cp
		p.onCharacter(start, p.lexer.I, p.state.lastIntValue)
		return true
	}

	if p.lexer.Eat(unicode_consts.ReverseSolidus) {
		if p.consumeClassEscape() {
			return true
		}

		p.raise("Invalid escape")

		p.lexer.Rewind(start)
	}

	return false
}

// ------------------------------------------------------------------------------
// ClassEscape ::
//
//	b
//	-
//	CharacterClassEscape
//	CharacterEscape
//
// https://tc39.es/ecma262/multipage/text-processing.html#prod-ClassEscape
// ------------------------------------------------------------------------------
func (p *Parser) consumeClassEscape() bool {
	start := p.lexer.I

	// b
	if p.lexer.Eat(unicode_consts.LatinSmallLetterB) {
		p.state.lastIntValue = unicode_consts.Backspace
		p.onCharacter(start-1, p.lexer.I, p.state.lastIntValue)
		return true
	}

	// -
	if p.lexer.Eat(unicode_consts.HyphenMinus) {
		p.state.lastIntValue = unicode_consts.HyphenMinus
		p.onCharacter(start-1, p.lexer.I, p.state.lastIntValue)
		return true
	}

	return p.consumeCharacterClass() || p.consumeCharacterEscape()
}

// ------------------------------------------------------------------------------
// CharacterClassEscape ::
//
//	d
//	D
//	s
//	S
//	w
//	W
//	p{ UnicodePropertyValueExpression }
//	P{ UnicodePropertyValueExpression }
//
// https://tc39.es/ecma262/multipage/text-processing.html#prod-CharacterClassEscape
// ------------------------------------------------------------------------------
func (p *Parser) consumeCharacterClassEscape() bool {
	return false
}

// ------------------------------------------------------------------------------
// CharacterEscape ::
//
//	ControlEscape
//	c AsciiLetter
//	0 [lookahead ∉ DecimalDigit]
//	HexEscapeSequence
//	RegExpUnicodeEscapeSequence
//	IdentityEscape
//
// https://tc39.es/ecma262/multipage/text-processing.html#prod-CharacterEscape
// ------------------------------------------------------------------------------
func (p *Parser) consumeCharacterEscape() bool {
	return false
}
