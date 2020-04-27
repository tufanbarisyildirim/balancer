package balancer

import (
	"sync/atomic"
	"time"
)

//Upstream is a mocking map between nodes and balancer
type Upstream struct {
	Node
	Weight           int
	load             int64
	requestCount     uint64
	totalRequestTime uint64
	nodeID           string
	healthy          bool
}

//IsHealthy returns if a node is healthy or not
func (u *Upstream) IsHealthy() bool {
	return u.healthy
}

//Load returns current load of a node
func (u *Upstream) Load() int64 {
	return atomic.LoadInt64(&u.load)
}

//IncreaseLoad increase load for a node
func (u *Upstream) IncreaseLoad() {
	atomic.AddInt64(&u.load, 1)
	atomic.AddUint64(&u.requestCount, 1)
}

//DecreaseLoad decrease load number of a node
func (u *Upstream) DecreaseLoad() {
	atomic.AddInt64(&u.load, -1)
}

//NodeID get name or address anything points to the
func (u *Upstream) NodeID() string {
	return u.nodeID
}

//DoRequest perform request mocking
func (u *Upstream) DoRequest() {
	defer u.DecreaseLoad()
	start := time.Now()
	u.IncreaseLoad()
	atomic.AddUint64(&u.totalRequestTime, uint64(time.Since(start)))
}

//IncreaseTime increase time for test
func (u *Upstream) IncreaseTime() {
	atomic.AddUint64(&u.totalRequestTime, uint64(time.Second*10))
}

//TotalRequestTime get sum of all request
func (u *Upstream) TotalRequestTime() uint64 {
	return atomic.LoadUint64(&u.totalRequestTime)
}

//TotalRequest total requests done to this node
func (u *Upstream) TotalRequest() uint64 {
	return atomic.LoadUint64(&u.requestCount)
}

//AverageResponseTime get overall average response time
//TODO(tufan): change it to get average in last x min?
func (u *Upstream) AverageResponseTime() time.Duration {
	if u.TotalRequest() == 0 {
		return 0
	}
	return time.Duration(u.TotalRequestTime() / u.TotalRequest())
}
