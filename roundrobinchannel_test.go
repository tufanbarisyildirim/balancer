package balancer

import (
	"reflect"
	"testing"
)

func TestRoundRobinChannel_SelectNode(t *testing.T) {
	type fields struct {
		robin uint32
	}
	type args struct {
		balancer *Balancer
		clientID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Node
	}{
		{
			name: "first and only comes",
			fields: fields{
				robin: 0, //to get first node
			},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: true,
							nodeID:  "127.0.0.2",
						},
					},
				},
			},
			want: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
			},
		},
		{
			name: "nil comes when no upstream",
			fields: fields{
				robin: 0, //to get first node
			},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{},
				},
			},
			want: nil,
		},
		{
			name: "nil comes when no healty upstream",
			fields: fields{
				robin: 0, //to get first node
			},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: false,
							nodeID:  "127.0.0.2",
						},
						&Upstream{
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
