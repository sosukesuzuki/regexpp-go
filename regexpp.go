package regexpp

func ParsePattern(source string, u bool) any {
	parser := NewParser(u)
	return parser.ParsePattern()
}
