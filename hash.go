package balancer

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

// findIndex finds consistent index using golang fast hash
// https://github.com/golang/go/blob/master/src/hash/fnv/fnv.go#L100 the code we previously  used
// has some allocations
// we need minimum footprint per deciding
func findIndex(s string, poolSize uint32) uint32 {
	var h32a uint32 = 2166136261
	for _, c := range []byte(s) {
		h32a *= 16777619
		h32a ^= uint32(c)
	}
	return h32a % poolSize
}
