package balancer

//Balancer select a node to send load
type Balancer struct {
	UpstreamPool []Node
	load         uint64
}
