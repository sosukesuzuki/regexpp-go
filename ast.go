package regexpp

type Node struct {
	Data N
	Loc  Loc
}

type Loc struct {
	start int
	end   int
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
