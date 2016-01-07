package Parser


type(
	Node struct {
		Id                           int
		Name, OtherName, Description string
		Parent                       *Node
	}
	DNode struct {
		Id          int `bson:"_id,omitempty" json:"_id"`
		Parent      int `bson:"parent,omitempty" json:"parent"`
		Name        string `bson:"name" json:"name"`
		OtherName   string `bson:"othername" json:"othername"`
		Description string `bson:"description" json:"description"`
	}

)

func (n *Node)ToDNode() *DNode {
	return &DNode{Id:n.Id, Name:n.Name, Parent: n.Parent.Id, OtherName:n.OtherName, Description:n.Description}
}
