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

func (n *NDisjunction) isNode() {}
func (n *NAlternative) isNode() {}

type NDisjunction struct {
	Left  Node
	Right Node
}

type NAlternative struct {
	Left  Node
	Right Node
}
