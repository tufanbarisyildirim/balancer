package balancer

//LeastTime select a node with fastest response time
type LeastTime struct{}

//SelectNode select next node has least response time in queue
func (lt *LeastTime) SelectNode(balancer *Balancer, clientID string) Node {
	var selectedNode Node

	for _, upstream := range balancer.UpstreamPool {
		if !upstream.IsHealthy() { //no one needs unhealth nodes.
			continue
		}

		if selectedNode == nil || selectedNode.GetAverageResponseTime() > upstream.GetAverageResponseTime() {
			selectedNode = upstream
			continue
		}
	}
	return selectedNode
}
