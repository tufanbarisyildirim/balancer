package balancer

//Selector Interface, policy to decide target
type SelectionPolicy interface {
	SelectNode(balancer *Balancer, clientID string) Node
}
