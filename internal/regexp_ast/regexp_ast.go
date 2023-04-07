package regexp_ast

type Loc struct {
	Start int
	End   int
}

type Node interface {
	isNode()
}

func (n *Pattern) isNode()     {}
func (n *Alternative) isNode() {}
func (n *Character) isNode()   {}

type Element interface {
	isElement()
}

func (n *Character) isElement() {}

type Pattern struct {
	Alternatives []Alternative
	Loc          Loc
}

type Alternative struct {
	Elements []Element
	Loc      Loc
}

type Character struct {
	value uint
	Loc   Loc
}
