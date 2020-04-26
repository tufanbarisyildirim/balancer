package balancer

import (
	"sync/atomic"
	"time"
)

//Upstream is a mocking map between nodes and balancer
type Upstream struct {
	Node
	Weight           int
	Load             int64
	RequestCount     uint64
	TotalRequestTime uint64
	Host             string
	Healthy          bool
}

//IsHealthy returns if a node is healthy or not
func (u *Upstream) IsHealthy() bool {
	return u.Healthy
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

//GetTotalRequestTime get sum of all request
func (u *Upstream) GetTotalRequestTime() uint64 {
	return atomic.LoadUint64(&u.TotalRequestTime)
}

//GetTotalRequest total requests done to this node
func (u *Upstream) GetTotalRequest() uint64 {
	return atomic.LoadUint64(&u.RequestCount)
}

//GetAverageResponseTime get overall average response time
//TODO(tufan): change it to get average in last x min?
func (u *Upstream) GetAverageResponseTime() time.Duration {
	return time.Duration(u.GetTotalRequestTime() / u.GetTotalRequest())
}
