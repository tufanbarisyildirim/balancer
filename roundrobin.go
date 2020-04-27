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
	poolSize := uint32(len(balancer.UpstreamPool))
	if poolSize == 0 {
		return nil
	}
	atomic.AddUint32(&r.RequestCount, 1)
	i := atomic.LoadUint32(&r.RequestCount)
	try := i + poolSize
	for ; i < try; i++ {
		upstream := balancer.UpstreamPool[i%poolSize]
		if upstream.IsHealthy() {
			return upstream
		}
		atomic.AddUint32(&r.RequestCount, 1)
	}
	return nil
}
