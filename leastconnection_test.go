package balancer

import (
	"reflect"
	"testing"
)

func TestLeastConnection_SelectNode(t *testing.T) {
	type args struct {
		balancer *Balancer
		clientID string
	}
	tests := []struct {
		name string
		lc   *LeastConnection
		args args
		want Node
	}{
		{
			name: "select the one and only we have",
			lc:   &LeastConnection{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: true,
							nodeID:  "127.0.0.1",
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.1",
			},
		},
		{
			name: "select the one that heav least connection",
			lc:   &LeastConnection{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: true,
							nodeID:  "127.0.0.1",
							load:    2,
						},
						&Upstream{
							healthy: true,
							nodeID:  "127.0.0.2",
							load:    1,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
		},
		{
			name: "select the one that heav least connection and healthy",
			lc:   &LeastConnection{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: true,
							nodeID:  "127.0.0.1",
							load:    2,
						},
						&Upstream{
							healthy: true,
							nodeID:  "127.0.0.2",
							load:    1,
						},
						&Upstream{
							healthy: false,
							nodeID:  "127.0.0.3",
							load:    0,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
		},
		{
			name: "return nil when no healthy node",
			lc:   &LeastConnection{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: false,
							nodeID:  "127.0.0.1",
							load:    2,
						},
						&Upstream{
							healthy: false,
							nodeID:  "127.0.0.2",
							load:    1,
						},
						&Upstream{
							healthy: false,
							nodeID:  "127.0.0.3",
							load:    0,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LeastConnection{}
			if got := lc.SelectNode(tt.args.balancer, tt.args.clientID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LeastConnection.SelectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
