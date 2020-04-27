package balancer

//Balancer select a node to send load
type Balancer struct {
	UpstreamPool []Node
	load         uint64
	Policy       SelectionPolicy
}

//Add add a node to balancer
func (b *Balancer) Add(node ...Node) {
	b.UpstreamPool = append(b.UpstreamPool, node...)
}

//Next select next available node
func (b *Balancer) Next(clientID string) Node {
	return b.Policy.SelectNode(b, clientID)
}
