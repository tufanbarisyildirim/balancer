package balancer

import (
	"sync/atomic"
	"time"
)

//MockNode is a mocking map between nodes and balancer
type MockNode struct {
	Node
	Weight           int
	load             int64
	requestCount     uint64
	totalRequestTime uint64
	nodeID           string
	healthy          bool
}

//IsHealthy returns if a node is healthy or not
func (mn *MockNode) IsHealthy() bool {
	return mn.healthy
}

//Load returns current load of a node
func (mn *MockNode) Load() int64 {
	return atomic.LoadInt64(&mn.load)
}

//IncreaseLoad increase load for a node
func (mn *MockNode) IncreaseLoad() {
	atomic.AddInt64(&mn.load, 1)
	atomic.AddUint64(&mn.requestCount, 1)
}

//DecreaseLoad decrease load number of a node
func (mn *MockNode) DecreaseLoad() {
	atomic.AddInt64(&mn.load, -1)
}

//NodeID get name or address anything points to the
func (mn *MockNode) NodeID() string {
	return mn.nodeID
}

//DoRequest perform request mocking
func (mn *MockNode) DoRequest() {
	defer mn.DecreaseLoad()
	start := time.Now()
	mn.IncreaseLoad()
	atomic.AddUint64(&mn.totalRequestTime, uint64(time.Since(start)))
}

//IncreaseTime increase time for test
func (mn *MockNode) IncreaseTime() {
	atomic.AddUint64(&mn.totalRequestTime, uint64(time.Second*10))
}

//TotalRequestTime get sum of all request
func (mn *MockNode) TotalRequestTime() uint64 {
	return atomic.LoadUint64(&mn.totalRequestTime)
}

//TotalRequest total requests done to this node
func (mn *MockNode) TotalRequest() uint64 {
	return atomic.LoadUint64(&mn.requestCount)
}

//AverageResponseTime get overall average response time
//TODO(tufan): change it to get average in last x min?
func (mn *MockNode) AverageResponseTime() time.Duration {
	if mn.TotalRequest() == 0 {
		return 0
	}
	return time.Duration(mn.TotalRequestTime() / mn.TotalRequest())
}
