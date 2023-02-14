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

type NDisjunction struct {
	elements []Node
}
