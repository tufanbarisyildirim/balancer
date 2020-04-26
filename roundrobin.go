package balancer

import (
	"sync/atomic"
)

//RoundRobin selection
type RoundRobin struct {
	RequestCount uint32
}

//SelectNode select next node in queue
func (r *RoundRobin) SelectNode(balancer *Balancer, clientID string) Node {
	poolSize := len(balancer.UpstreamPool)
	
	if poolSize == 0 {
		return nil
	}
	for i := 0; i < poolSize; i++ {
		atomic.AddUint32(&r.RequestCount, 1)
		if upstream := balancer.UpstreamPool[atomic.LoadUint32(&r.RequestCount)%uint32(poolSize)]; upstream.IsHealthy() {
			return upstream
		}
	}
	return nil
}
