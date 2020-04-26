package balancer

//LeastConnection select a node with lower load
type LeastConnection struct{}

//SelectNode select next node in queue
func (lc *LeastConnection) SelectNode(balancer *Balancer, clientID string) Node {
	var selectedNode Node

	for _, upstream := range balancer.UpstreamPool {
		if !upstream.IsHealthy() { //no one needs unhealth nodes.
			continue
		}

		if selectedNode == nil || selectedNode.GetLoad() > upstream.GetLoad() {
			selectedNode = upstream
			continue
		}
	}
	return selectedNode
}
