package balancer

import (
	"time"
)

//RoundRobinChannel implements roundrobin using go channels
type RoundRobinChannel struct {
	pool chan Node
}

//SelectNode select next node in queue
func (rc *RoundRobinChannel) SelectNode(balancer *Balancer, clientID string) Node {
	poolSize := len(balancer.UpstreamPool)
	if poolSize == 0 {
		return nil
	}

	if rc.pool == nil {
		rc.pool = make(chan Node, poolSize)
	}

	if len(rc.pool) == 0 {
		go func() {
			for i := 0; i < poolSize; i++ {
				rc.pool <- balancer.UpstreamPool[i]
			}
		}()
	}

	try := 0
	for {
		select {
		case n := <-rc.pool:
			go func() {
				rc.pool <- n
			}()
			if n.IsHealthy() {
				return n
			}
		case <-time.After(time.Millisecond * 2): //empty queue? wait 2ms
			goto goon
		}

	goon:
		try = try + 1
		if try >= poolSize {
			return nil
		}
	}
}
