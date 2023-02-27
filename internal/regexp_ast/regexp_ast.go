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
	isNode()
}

func (n *NPattern) isNode()     {}
func (n *NDisjunction) isNode() {}
func (n *NAlternative) isNode() {}
func (n *NAtom) isNode()        {}

type NPattern struct {
	Value *Node
}

type NDisjunction struct {
	Left  *Node
	Right *Node
}

type NAlternative struct {
	Left  *Node
	Right *Node
}

type NAtom struct {
	Value *Node
}
