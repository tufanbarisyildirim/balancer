package balancer

import "time"

//Node a node that handles load
type Node interface {
	IsHealthy() bool
	GetTotalRequest() uint64
	GetAverageResponseTime() time.Duration
	GetLoad() int64 
	GetHost() string
}
