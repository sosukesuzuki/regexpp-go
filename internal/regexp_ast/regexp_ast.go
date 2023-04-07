package regexp_ast

type Node struct {
	Data N
	Loc  Loc
}

type Loc struct {
	Start int
	End   int
}

type N interface {
	isN()
}

func (n *Pattern) isN()     {}
func (n *Alternative) isN() {}
func (n *Character) isN()   {}

type Element interface {
	isElement()
}

func (n *Character) isElement() {}

type Pattern struct {
	Alternatives []Alternative
}

type Alternative struct {
	Elements []Element
}

type Character struct {
	value uint
}
