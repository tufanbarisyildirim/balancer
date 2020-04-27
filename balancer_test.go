package balancer

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestBalancer_Add(t *testing.T) {
	type fields struct {
		UpstreamPool []Node
		load         uint64
		selector     SelectionPolicy
	}
	type args struct {
		node []Node
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "add 1 node",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.1:8000",
					},
				},
			},
			want: 1,
		},
		{
			name: "add 2 node",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.1:8000",
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.2:8000",
					},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				UpstreamPool: tt.fields.UpstreamPool,
				load:         tt.fields.load,
				Policy:     tt.fields.selector,
			}
			b.Add(tt.args.node...)

			if got := len(b.UpstreamPool); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBalancer_Next(t *testing.T) {
	type fields struct {
		UpstreamPool []Node
		load         uint64
		selector     SelectionPolicy
	}
	type args struct {
		clientID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want1  Node
		want2  Node
		want3  Node
		want4  Node
		want5  Node
	}{
		{
			name: "roundrobin next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						Healthy: false,
						Host:    "127.0.0.1",
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.2",
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.3",
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.4",
					},
				},
				load:     0,
				selector: SelectionPolicy(&RoundRobin{}),
			},
			want1: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
			},
			want2: &Upstream{
				Healthy: true,
				Host:    "127.0.0.3",
			},
			want3: &Upstream{
				Healthy: true,
				Host:    "127.0.0.4",
			},
			want4: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
			},
			want5: &Upstream{
				Healthy: true,
				Host:    "127.0.0.3",
			},
		},
		{
			name: "leastconnection next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						Healthy: false,
						Host:    "127.0.0.1",
						Load:    0,
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.2",
						Load:    1,
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.3",
						Load:    2,
					},
					&Upstream{
						Healthy: true,
						Host:    "127.0.0.4",
						Load:    3,
					},
				},
				load:     0,
				selector: SelectionPolicy(&LeastConnection{}),
			},
			want1: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
			want2: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
			want3: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
			want4: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
			want5: &Upstream{
				Healthy: true,
				Host:    "127.0.0.2",
				Load:    1,
			},
		},
		{
			name: "leasttimeconnection next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						Healthy: false,
						Host:    "127.0.0.1",
						Load:    0,
					},
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.2",
						Load:             1,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.3",
						Load:             2,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.4",
						Load:             3,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 5),
					},
				},
				load:     0,
				selector: SelectionPolicy(&LeastTime{}),
			},
			want1: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want2: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want3: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want4: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want5: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
		},
		{
			name: "hash next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.4",
						Load:             3,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 5),
					},
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.3",
						Load:             2,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						Healthy:          false,
						Host:             "127.0.0.1",
						Load:             0,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.2",
						Load:             1,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						Healthy:          true,
						Host:             "127.0.0.4",
						Load:             3,
						RequestCount:     5000,
						TotalRequestTime: uint64(time.Second * 5),
					},
				},
				load:     0,
				selector: SelectionPolicy(&Hash{}),
			},
			want1: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want2: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want3: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want4: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
			want5: &Upstream{
				Healthy:          true,
				Host:             "127.0.0.4",
				Load:             3,
				RequestCount:     5000,
				TotalRequestTime: uint64(time.Second * 5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				UpstreamPool: tt.fields.UpstreamPool,
				load:         tt.fields.load,
				Policy:     tt.fields.selector,
			}
			if got := b.Next(tt.args.clientID); !reflect.DeepEqual(got, tt.want1) {
				t.Errorf("Balancer.Next() = %v, want1 %v", got, tt.want1)
			}
			if got := b.Next(tt.args.clientID); !reflect.DeepEqual(got, tt.want2) {
				t.Errorf("Balancer.Next() = %v, want2 %v", got, tt.want2)
			}
			if got := b.Next(tt.args.clientID); !reflect.DeepEqual(got, tt.want3) {
				t.Errorf("Balancer.Next() = %v, want3 %v", got, tt.want3)
			}
			if got := b.Next(tt.args.clientID); !reflect.DeepEqual(got, tt.want4) {
				t.Errorf("Balancer.Next() = %v, want4 %v", got, tt.want4)
			}
			if got := b.Next(tt.args.clientID); !reflect.DeepEqual(got, tt.want5) {
				t.Errorf("Balancer.Next() = %v, want5 %v", got, tt.want5)
			}
		})
	}
}

func genNodes() []Node {
	nodes := make([]Node, 0)
	for i := 0; i < 10; i++ {
		nodes = append(nodes, &Upstream{
			Weight:           1,
			Host:             fmt.Sprintf("127.0.0.%d:8080", i),
			Healthy:          i%2 == 0,
			RequestCount:     uint64(time.Now().Nanosecond()),
			TotalRequestTime: uint64(time.Now().Nanosecond() * 300),
			Load:             int64(time.Now().Second()),
		})
	}
	return nodes
}

//BenchmarkNextRoundRobin benchmark roundrobin algorithm
func BenchmarkNextRoundRobin(b *testing.B) {

	balancer := Balancer{
		UpstreamPool: genNodes(),
		load:         0,
		Policy:     &RoundRobin{},
	}

	for n := 0; n < b.N; n++ {
		upstream := balancer.Next("127.0.0.1").(*Upstream)
		if n%50 == 0 {
			upstream.IncreaseLoad()
		}
	}
}

//BenchmarkNextHash consistant hashing
func BenchmarkNextHash(b *testing.B) {

	balancer := Balancer{
		UpstreamPool: genNodes(),
		load:         0,
		Policy:     &Hash{},
	}

	for n := 0; n < b.N; n++ {
		upstream := balancer.Next("127.0.0.1").(*Upstream)
		if n%50 == 0 {
			upstream.IncreaseLoad()
		}
	}
}

//BenchmarkNextLeastConnection least connection
func BenchmarkNextLeastConnection(b *testing.B) {

	balancer := Balancer{
		UpstreamPool: genNodes(),
		load:         0,
		Policy:     &LeastConnection{},
	}

	for n := 0; n < b.N; n++ {
		upstream := balancer.Next("127.0.0.1").(*Upstream)
		if n%50 == 0 {
			upstream.IncreaseLoad()
		}
	}
}

//BenchmarkNextLeastTime
func BenchmarkNextLeastTime(b *testing.B) {

	balancer := Balancer{
		UpstreamPool: genNodes(),
		load:         0,
		Policy:     &LeastTime{},
	}

	for n := 0; n < b.N; n++ {
		upstream := balancer.Next("127.0.0.1").(*Upstream)
		if n%5 == 0 {
			upstream.IncreaseTime()
		}
	}
}
