package balancer

import (
	"reflect"
	"testing"
	"time"
)

func TestLeastTime_SelectNode(t *testing.T) {
	type args struct {
		balancer *Balancer
		clientID string
	}
	tests := []struct {
		name string
		lt   *LeastTime
		args args
		want Node
	}{
		{
			name: "select the one and only we have",
			lt:   &LeastTime{},
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
			name: "select the one that has least connection response time",
			lt:   &LeastTime{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							Healthy:          true,
							Host:             "127.0.0.1",
							Load:             2,
							RequestCount:     1000,
							TotalRequestTime: uint64(time.Second * 10),
						},
						&Upstream{
							Healthy:          true,
							Host:             "127.0.0.2",
							Load:             1,
							RequestCount:     1000,
							TotalRequestTime: uint64(time.Second * 5),
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.2",
				Load:             1,
				RequestCount:     1000,
				TotalRequestTime: uint64(time.Second * 5),
			},
		},
		{
			name: "select the one that have least connection time and healthy",
			lt:   &LeastTime{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							Healthy:          true,
							Host:             "127.0.0.1",
							Load:             2,
							RequestCount:     5000,
							TotalRequestTime: uint64(time.Second * 5),
						},
						&Upstream{
							Healthy:          true,
							Host:             "127.0.0.2",
							Load:             1,
							RequestCount:     5000,
							TotalRequestTime: uint64(time.Second * 4),
						},
						&Upstream{
							Healthy:          false,
							Host:             "127.0.0.3",
							Load:             0,
							RequestCount:     5000,
							TotalRequestTime: 100,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.2",
				Load:             1,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 4),
			},
		},
		{
			name: "return nill cuz we have no healthy one",
			lt:   &LeastTime{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&Upstream{
							Healthy:          false,
							Host:             "127.0.0.1",
							Load:             2,
							RequestCount:     5000,
							TotalRequestTime: uint64(time.Second * 5),
						},
						&Upstream{
							Healthy:          false,
							Host:             "127.0.0.2",
							Load:             1,
							RequestCount:     5000,
							TotalRequestTime: uint64(time.Second * 4),
						},
						&Upstream{
							Healthy:          false,
							Host:             "127.0.0.3",
							Load:             0,
							RequestCount:     5000,
							TotalRequestTime: 100,
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
			lt := &LeastTime{}
			if got := lt.SelectNode(tt.args.balancer, tt.args.clientID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LeastTime.SelectNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
