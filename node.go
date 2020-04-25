package balancer

import (
	"sync/atomic"
	"time"
)

//Node a node that handles load
type Node interface {
	IsHealthy() bool
	GetTotalRequest() uint64
	GetLoad() int64
	IncreaseLoad()
	DecreaseLoad()
	GetHost() string
}

//Pool node list
type Pool []Node

//Upstream is a mocking map between nodes and balancer
type Upstream struct {
	Node
	Weight           int
	Load             int64
	RequestCount     uint64
	TotalRequestTime uint64
	Host             string
}

//IsHealthy returns if a node is healthy or not
func (u *Upstream) IsHealthy() bool {
	return true
}

//GetLoad returns current load of a node
func (u *Upstream) GetLoad() int64 {
	return atomic.LoadInt64(&u.Load)
}

//IncreaseLoad increase load for a node
func (u *Upstream) IncreaseLoad() {
	atomic.AddInt64(&u.Load, 1)
	atomic.AddUint64(&u.RequestCount, 1)
}

//DecreaseLoad decrease load number of a node
func (u *Upstream) DecreaseLoad() {
	atomic.AddInt64(&u.Load, -1)
}

//GetHost get name or address anything points to the
func (u *Upstream) GetHost() string {
	return u.Host
}

//DoRequest perform request mocking
func (u *Upstream) DoRequest() {
	defer u.DecreaseLoad()
	start := time.Now()
	u.IncreaseLoad()
	atomic.AddUint64(&u.TotalRequestTime, uint64(time.Since(start)))
}

//GetTotalRequest total requests done to this node
func (u *Upstream) GetTotalRequest() uint64 {
	return atomic.LoadUint64(&u.RequestCount)
}
