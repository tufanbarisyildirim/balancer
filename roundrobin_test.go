package balancer

import (
	"reflect"
	"testing"
)

func TestRoundRobin_SelectNode(t *testing.T) {
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
							Healthy: true,
							Host:    "127.0.0.2",
						},
					},
				},
			},
			want: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
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
							Healthy: false,
							Host:    "127.0.0.2",
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoundRobin{
				RequestCount: tt.fields.robin,
			}
			if got := r.SelectNode(tt.args.balancer, tt.args.clientID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundRobin.SelectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
