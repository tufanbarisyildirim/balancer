package balancer

//Selector Interface, policy to decide target
type Selector interface {
	SelectNode(Pool) *Node
}
