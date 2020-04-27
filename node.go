package balancer

import "time"

//Node a node that handles load
type Node interface {
	IsHealthy() bool
	TotalRequest() uint64
	AverageResponseTime() time.Duration
	Load() int64
	NodeID() string
}
