package balancer

import (
	"testing"
)

func TestUpstream_IsHealthy(t *testing.T) {
	tests := []struct {
		name     string
		upstream MockNode
		want     bool
	}{
		{
			name: "healty trivial test",
			upstream: MockNode{
				healthy: true,
			},
			want: true,
		},
		{
			name: "healty trivial test",
			upstream: MockNode{
				healthy: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.upstream
			if got := u.IsHealthy(); got != tt.want {
				t.Errorf("MockNode.IsHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpstream_IncreaseLoad(t *testing.T) {
	tests := []struct {
		name             string
		upstream         *MockNode
		call             func(u *MockNode)
		wantRequestCount uint64
		wantLoad         int64
		wantHost         string
	}{
		{
			name: "IncreaseLoad 1 time",
			upstream: &MockNode{
				nodeID: "127.0.0.1",
			},
			call: func(u *MockNode) {
				u.IncreaseLoad()
			},
			wantRequestCount: uint64(1),
			wantLoad:         int64(1),
			wantHost:         "127.0.0.1",
		},
		{
			name: "IncreaseLoad 2 time, decrease 1 time",
			upstream: &MockNode{
				nodeID: "127.0.0.1",
			},
			call: func(u *MockNode) {
				u.IncreaseLoad()
				u.IncreaseLoad()
				u.DecreaseLoad()
			},
			wantRequestCount: uint64(2),
			wantLoad:         int64(1),
			wantHost:         "127.0.0.1",
		},

		{
			name: "Do requests, no load but increased (done) request",
			upstream: &MockNode{
				nodeID: "127.0.0.1",
			},
			call: func(u *MockNode) {
				u.DoRequest()
				u.DoRequest()
				u.DoRequest()
				u.DoRequest()
			},
			wantRequestCount: uint64(4),
			wantLoad:         int64(0),
			wantHost:         "127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.upstream
			tt.call(u)
			if tt.wantLoad != u.Load() {
				t.Errorf("MockNode.Load() = %v, want %v", u.Load(), tt.wantLoad)
			}
			if tt.wantRequestCount != u.TotalRequest() {
				t.Errorf("MockNode.TotalRequest() = %v, want %v", u.TotalRequest(), tt.wantRequestCount)
			}
			if tt.wantHost != u.NodeID() {
				t.Errorf("MockNode.NodeID() = %v, want %v", u.NodeID(), tt.wantHost)
			}
		})
	}
}
