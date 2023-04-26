package regexp_ast

type Loc struct {
	Start int
	End   int
}

type Node interface {
	isNode()
	GetParent() Node
	SetEnd(end int)
}

func (n *Pattern) isNode()        {}
func (n *Alternative) isNode()    {}
func (n *Character) isNode()      {}
func (n *CharacterClass) isNode() {}

func (n *Pattern) GetParent() Node        { return nil }
func (n *Alternative) GetParent() Node    { return n.Parent }
func (n *Character) GetParent() Node      { return n.Parent }
func (n *CharacterClass) GetParent() Node { return n.Parent }

func (n *Pattern) SetEnd(end int)        { n.Loc.End = end }
func (n *Alternative) SetEnd(end int)    { n.Loc.End = end }
func (n *Character) SetEnd(end int)      { n.Loc.End = end }
func (n *CharacterClass) SetEnd(end int) { n.Loc.End = end }

type Element interface {
	isElement()
}

func (n *Character) isElement() {}

type Pattern struct {
	Alternatives []*Alternative
	Loc          Loc
}

type Alternative struct {
	Elements []Element
	Parent   Node
	Loc      Loc
}

type Character struct {
	Value  int
	Parent Node
	Loc    Loc
}

type CharacterClass struct {
	Parent   Node
	Loc      Loc
	Negate   bool
	Elements []Element
}
