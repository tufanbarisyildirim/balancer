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
			name:     "healty trivial test",
			upstream: Upstream{},
			want:     true,
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
				Host: "127.0.0.1",
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
				Host: "127.0.0.1",
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
				Host: "127.0.0.1",
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
			if tt.wantLoad != u.GetLoad() {
				t.Errorf("Upstream.GetLoad() = %v, want %v", u.GetLoad(), tt.wantLoad)
			}
			if tt.wantRequestCount != u.GetTotalRequest() {
				t.Errorf("Upstream.GetTotalRequest() = %v, want %v", u.GetTotalRequest(), tt.wantRequestCount)
			}
			if tt.wantHost != u.GetHost() {
				t.Errorf("Upstream.GetHost() = %v, want %v", u.GetHost(), tt.wantHost)
			}
		})
	}
}
