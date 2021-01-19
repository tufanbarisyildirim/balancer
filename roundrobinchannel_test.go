package balancer

import (
	"reflect"
	"testing"
)

func TestRoundRobinChannel_SelectNode(t *testing.T) {
	type args struct {
		balancer *Balancer
		clientID string
	}
	tests := []struct {
		name   string
		args   args
		want   Node
	}{
		{
			name: "first and only comes",
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&MockNode{
							healthy: true,
							nodeID:  "127.0.0.2",
						},
					},
				},
			},
			want: &MockNode{
				healthy: true,
				nodeID:  "127.0.0.2",
			},
		},
		{
			name: "nil comes when no upstream",
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{},
				},
			},
			want: nil,
		},
		{
			name: "nil comes when no healthy upstream",
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&MockNode{
							healthy: false,
							nodeID:  "127.0.0.2",
						},
						&MockNode{
							healthy: false,
							nodeID:  "127.0.0.1",
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoundRobinChannel{}
			if got := r.SelectNode(tt.args.balancer, tt.args.clientID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundRobin.SelectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
