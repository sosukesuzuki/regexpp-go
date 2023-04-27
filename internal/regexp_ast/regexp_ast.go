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
func (n *AnyCharacterSet) isElement() {}

type CharacterSet interface {
	isCharacterSet()
}

func (n *AnyCharacterSet) isCharacterSet() {}

type Pattern struct {
	Loc          Loc
	Alternatives []*Alternative
}

type Alternative struct {
	Parent   Node `json:"-"`
	Elements []Element
	Loc      Loc
}

type Character struct {
	Parent Node `json:"-"`
	Loc    Loc
	Value  int
}

type CharacterClass struct {
	Parent   Node `json:"-"`
	Loc      Loc
	Negate   bool
	Elements []Element
}

// Dot
type AnyCharacterSet struct {
	Parent Node `json:"-"`
	Loc Loc
}
