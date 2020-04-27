package balancer

import (
	"hash/fnv"
)

//Hash jump to a node using consistent hash
type Hash struct {
}

//SelectNode select a node based on hash
func (h *Hash) SelectNode(balancer *Balancer, clientID string) Node {
	poolSize := uint32(len(balancer.UpstreamPool))
	if poolSize == 0 {
		return nil
	}
	index := findIndex(clientID, poolSize)
	if upstream := balancer.UpstreamPool[index]; upstream.IsHealthy() {
		return upstream
	}
	for i := uint32(0); i < poolSize; i++ {
		if upstream := balancer.UpstreamPool[i]; upstream.IsHealthy() {
			return upstream
		}
	}
	return nil
}

// findIndex finds consistant index using golang fast hash
func findIndex(s string, poolSize uint32) uint32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(s))
	return hasher.Sum32() % poolSize
}
