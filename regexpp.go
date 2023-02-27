package regexpp

import "github.com/sosukesuzuki/regexpp-go/internal/parser"

func ParsePattern(source string, u bool) any {
	parser := parser.NewParser(source, u)
	return parser.ParsePattern()
}
