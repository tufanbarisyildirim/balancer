package balancer

//Node a node that handles load
type Node interface {
	IsHealthy() bool
	GetTotalRequest() uint64
	GetLoad() int64
	IncreaseLoad()
	DecreaseLoad()
	GetHost() string
}
