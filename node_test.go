package balancer

import (
	"testing"
)

func TestUpstream_IsHealthy(t *testing.T) {
	tests := []struct {
		name     string
		upstream Upstream
		want     bool
	}{
		{
			name: "healty trivial test",
			upstream: Upstream{
				healthy: true,
			},
			want: true,
		},
		{
			name: "healty trivial test",
			upstream: Upstream{
				healthy: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := tt.upstream
			if got := u.IsHealthy(); got != tt.want {
				t.Errorf("Upstream.IsHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpstream_IncreaseLoad(t *testing.T) {
	tests := []struct {
		name             string
		upstream         *Upstream
		call             func(u *Upstream)
		wantRequestCount uint64
		wantLoad         int64
		wantHost         string
	}{
		{
			name: "IncreaseLoad 1 time",
			upstream: &Upstream{
				nodeID: "127.0.0.1",
			},
			call: func(u *Upstream) {
				u.IncreaseLoad()
			},
			wantRequestCount: uint64(1),
			wantLoad:         int64(1),
			wantHost:         "127.0.0.1",
		},
		{
			name: "IncreaseLoad 2 time, decrease 1 time",
			upstream: &Upstream{
				nodeID: "127.0.0.1",
			},
			call: func(u *Upstream) {
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
			upstream: &Upstream{
				nodeID: "127.0.0.1",
			},
			call: func(u *Upstream) {
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
				t.Errorf("Upstream.Load() = %v, want %v", u.Load(), tt.wantLoad)
			}
			if tt.wantRequestCount != u.TotalRequest() {
				t.Errorf("Upstream.TotalRequest() = %v, want %v", u.TotalRequest(), tt.wantRequestCount)
			}
			if tt.wantHost != u.NodeID() {
				t.Errorf("Upstream.NodeID() = %v, want %v", u.NodeID(), tt.wantHost)
			}
		})
	}
}
