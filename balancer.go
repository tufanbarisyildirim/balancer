package balancer

//Balancer select a node to send load
type Balancer struct {
	UpstreamPool Pool
	load         uint64
}
