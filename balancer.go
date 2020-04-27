package balancer

import "sync"

//Balancer select a node to send load
type Balancer struct {
	UpstreamPool []Node
	load         uint64
	Policy       SelectionPolicy
	m            *sync.RWMutex
}

//NewBalancer create a new balancer with default properties
func NewBalancer() *Balancer {
	return &Balancer{
		UpstreamPool: make([]Node, 0),
		load:         0,
		Policy:       &RoundRobin{},
		m:            &sync.RWMutex{},
	}
}

//Add add a node to balancer
func (b *Balancer) Add(node ...Node) {
	b.UpstreamPool = append(b.UpstreamPool, node...)
}

//Remove remove a node from balancer
func (b *Balancer) Remove(nodeID string) {
	b.m.Lock()
	defer b.m.Unlock()
	i := 0
	for _, upstream := range b.UpstreamPool {
		if upstream.NodeID() != nodeID {
			b.UpstreamPool[i] = upstream
			i++
		}
	}
	b.UpstreamPool = b.UpstreamPool[:i]
}

//Next select next available node
func (b *Balancer) Next(clientID string) Node {
	defer b.m.Unlock()
	b.m.Lock()
	return b.Policy.SelectNode(b, clientID)

}
