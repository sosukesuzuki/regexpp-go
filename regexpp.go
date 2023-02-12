package regexpp

func ParsePattern(source string, u bool) any {
	parser := NewParser(source, u)
	return parser.ParsePattern()
}
