package balancer

import jump "github.com/lithammer/go-jump-consistent-hash"

//Hash jump to a node using consistent hash
type Hash struct {
}

//SelectNode select a node based on hash
func (h *Hash) SelectNode(balancer *Balancer, clientID string) Node {
	poolSize := int32(len(balancer.UpstreamPool))
	if poolSize == 0 {
		return nil
	}
	index := jump.HashString(clientID, poolSize, jump.NewCRC64()) % poolSize
	for i := int32(0); i < poolSize; i++ {
		index += i
		if upstream := balancer.UpstreamPool[index%poolSize]; upstream.IsHealthy() {
			return upstream
		}
	}
	return nil
}
