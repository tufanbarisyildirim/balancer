package balancer

import (
	"reflect"
	"testing"
)

func TestHash_SelectNode(t *testing.T) {
	type args struct {
		balancer *Balancer
		clientID string
	}
	tests := []struct {
		name string
		h    *Hash
		args args
		want Node
	}{
		{
			name: "same node when we have only 1 node",
			h:    &Hash{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							nodeID:  "127.0.0.1",
							healthy: true,
						},
					},
				},
			},
			want: &Upstream{
				nodeID:  "127.0.0.1",
				healthy: true,
			},
		},
		{
			name: "empty upstreams should return nil",
			h:    &Hash{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{},
				},
			},
			want: nil,
		},
		{
			name: "no healty upstream? then nil",
			h:    &Hash{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							healthy: false,
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.SelectNode(tt.args.balancer, tt.args.clientID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hash.SelectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
