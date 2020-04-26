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
							Healthy: true,
							Host:    "127.0.0.1",
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				Healthy: true,
				Host:    "127.0.0.1",
			},
		},
		{
			name: "select the one that heav least connection",
			lc:   &LeastConnection{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							Healthy: true,
							Host:    "127.0.0.1",
							Load:    2,
						},
						&Upstream{
							Healthy: true,
							Host:    "127.0.0.2",
							Load:    1,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
		},
		{
			name: "select the one that heav least connection and healthy",
			lc:   &LeastConnection{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							Healthy: true,
							Host:    "127.0.0.1",
							Load:    2,
						},
						&Upstream{
							Healthy: true,
							Host:    "127.0.0.2",
							Load:    1,
						},
						&Upstream{
							Healthy: false,
							Host:    "127.0.0.3",
							Load:    0,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
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
