package regexpp

import "github.com/sosukesuzuki/regexpp-go/internal/parser"

func ParsePattern(source string, u bool) {
	parser := parser.NewParser(source, u)
	 parser.ParsePattern()
}
