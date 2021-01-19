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
						&MockNode{
							healthy: true,
							nodeID:  "127.0.0.1",
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &MockNode{
				healthy: true,
				nodeID:  "127.0.0.1",
			},
		},
		{
			name: "select the one that has least connection response time",
			lt:   &LeastTime{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&MockNode{
							healthy:          true,
							nodeID:           "127.0.0.1",
							load:             2,
							requestCount:     1000,
							totalRequestTime: uint64(time.Second * 10),
						},
						&MockNode{
							healthy:          true,
							nodeID:           "127.0.0.2",
							load:             1,
							requestCount:     1000,
							totalRequestTime: uint64(time.Second * 5),
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &MockNode{
				healthy:          true,
				nodeID:           "127.0.0.2",
				load:             1,
				requestCount:     1000,
				totalRequestTime: uint64(time.Second * 5),
			},
		},
		{
			name: "select the one that have least connection time and healthy",
			lt:   &LeastTime{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&MockNode{
							healthy:          true,
							nodeID:           "127.0.0.1",
							load:             2,
							requestCount:     5000,
							totalRequestTime: uint64(time.Second * 5),
						},
						&MockNode{
							healthy:          true,
							nodeID:           "127.0.0.2",
							load:             1,
							requestCount:     5000,
							totalRequestTime: uint64(time.Second * 4),
						},
						&MockNode{
							healthy:          false,
							nodeID:           "127.0.0.3",
							load:             0,
							requestCount:     5000,
							totalRequestTime: 100,
						},
					},
				},
				clientID: "127.0.0.1",
			},
			want: &MockNode{
				healthy:          true,
				nodeID:           "127.0.0.2",
				load:             1,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 4),
			},
		},
		{
			name: "return nill cuz we have no healthy one",
			lt:   &LeastTime{},
			args: args{
				balancer: &Balancer{
					UpstreamPool: []Node{
						&MockNode{
							healthy:          false,
							nodeID:           "127.0.0.1",
							load:             2,
							requestCount:     5000,
							totalRequestTime: uint64(time.Second * 5),
						},
						&MockNode{
							healthy:          false,
							nodeID:           "127.0.0.2",
							load:             1,
							requestCount:     5000,
							totalRequestTime: uint64(time.Second * 4),
						},
						&MockNode{
							healthy:          false,
							nodeID:           "127.0.0.3",
							load:             0,
							requestCount:     5000,
							totalRequestTime: 100,
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
