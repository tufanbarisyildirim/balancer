package balancer

import (
	"fmt"
	"reflect"
	"sync"
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
						healthy: true,
						nodeID:  "127.0.0.1:8000",
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
						healthy: true,
						nodeID:  "127.0.0.1:8000",
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.2:8000",
					},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBalancer()
			b.UpstreamPool = tt.fields.UpstreamPool
			b.load = tt.fields.load
			b.Policy = tt.fields.selector

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
						healthy: false,
						nodeID:  "127.0.0.1",
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.2",
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.3",
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.4",
					},
				},
				load:     0,
				selector: SelectionPolicy(&RoundRobin{}),
			},
			want1: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
			},
			want2: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.3",
			},
			want3: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.4",
			},
			want4: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
			},
			want5: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.3",
			},
		},
		{
			name: "leastconnection next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						healthy: false,
						nodeID:  "127.0.0.1",
						load:    0,
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.2",
						load:    1,
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.3",
						load:    2,
					},
					&Upstream{
						healthy: true,
						nodeID:  "127.0.0.4",
						load:    3,
					},
				},
				load:     0,
				selector: SelectionPolicy(&LeastConnection{}),
			},
			want1: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
			want2: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
			want3: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
			want4: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
			want5: &Upstream{
				healthy: true,
				nodeID:  "127.0.0.2",
				load:    1,
			},
		},
		{
			name: "leasttimeconnection next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						healthy: false,
						nodeID:  "127.0.0.1",
						load:    0,
					},
					&Upstream{
						healthy:          true,
						nodeID:           "127.0.0.2",
						load:             1,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						healthy:          true,
						nodeID:           "127.0.0.3",
						load:             2,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						healthy:          true,
						nodeID:           "127.0.0.4",
						load:             3,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 5),
					},
				},
				load:     0,
				selector: SelectionPolicy(&LeastTime{}),
			},
			want1: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want2: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want3: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want4: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want5: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
		},
		{
			name: "leastconnection empty pool test",
			fields: fields{
				UpstreamPool: []Node{},
				load:         0,
				selector:     SelectionPolicy(&LeastConnection{}),
			},
			want1: nil,
			want2: nil,
			want3: nil,
			want4: nil,
			want5: nil,
		},
		{
			name: "leasttimeconnection empty pool test",
			fields: fields{
				UpstreamPool: []Node{},
				load:         0,
				selector:     SelectionPolicy(&LeastTime{}),
			},
			want1: nil,
			want2: nil,
			want3: nil,
			want4: nil,
			want5: nil,
		},
		{
			name: "hash next test",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{
						healthy:          true,
						nodeID:           "127.0.0.4",
						load:             3,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 5),
					},
					&Upstream{
						healthy:          false,
						nodeID:           "127.0.0.3",
						load:             2,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						healthy:          false,
						nodeID:           "127.0.0.1",
						load:             0,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						healthy:          true,
						nodeID:           "127.0.0.2",
						load:             1,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 10),
					},
					&Upstream{
						healthy:          true,
						nodeID:           "127.0.0.4",
						load:             3,
						requestCount:     5000,
						totalRequestTime: uint64(time.Second * 5),
					},
				},
				load:     0,
				selector: SelectionPolicy(&Hash{}),
			},
			want1: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want2: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want3: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want4: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
			want5: &Upstream{
				healthy:          true,
				nodeID:           "127.0.0.4",
				load:             3,
				requestCount:     5000,
				totalRequestTime: uint64(time.Second * 5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBalancer()

			b.UpstreamPool = tt.fields.UpstreamPool
			b.load = tt.fields.load
			b.Policy = tt.fields.selector

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
			nodeID:           fmt.Sprintf("127.0.0.%d:8080", i),
			healthy:          i%2 == 0,
			requestCount:     uint64(time.Now().Nanosecond()),
			totalRequestTime: uint64(time.Now().Nanosecond() * 300),
			load:             int64(time.Now().Second()),
		})
	}
	return nodes
}

//BenchmarkNextRoundRobin benchmark roundrobin algorithm
func BenchmarkNextRoundRobin(b *testing.B) {

	balancer := Balancer{
		UpstreamPool: genNodes(),
		load:         0,
		Policy:       &RoundRobin{},
		m:            &sync.RWMutex{},
	}

	for n := 0; n < b.N; n++ {
		upstream := balancer.Next("127.0.0.1").(*Upstream)
		if n%50 == 0 {
			upstream.IncreaseLoad()
		}
	}
}


//BenchmarkNextRoundRobinChannel benchmark roundrobin algorithm using go channels
func BenchmarkNextRoundRobinChannel(b *testing.B) {

	balancer := Balancer{
		UpstreamPool: genNodes(),
		load:         0,
		Policy:       &RoundRobinChannel{},
		m:            &sync.RWMutex{},
	}

	for n := 0; n < b.N; n++ {
		us:=balancer.Next("127.0.0.1")
		if us == nil{
			continue
		}
		upstream := us.(*Upstream)
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
		Policy:       &Hash{},
		m:            &sync.RWMutex{},
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
		Policy:       &LeastConnection{},
		m:            &sync.RWMutex{},
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
		Policy:       &LeastTime{},
		m:            &sync.RWMutex{},
	}

	for n := 0; n < b.N; n++ {
		upstream := balancer.Next("127.0.0.1").(*Upstream)
		if n%5 == 0 {
			upstream.IncreaseTime()
		}
	}
}

func TestBalancer_Remove(t *testing.T) {
	type fields struct {
		UpstreamPool []Node
		load         uint64
		Policy       SelectionPolicy
		m            *sync.RWMutex
	}
	type args struct {
		nodeID string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
	}{
		{
			name: "remove one",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{nodeID: "1"},
					&Upstream{nodeID: "2"},
				},
				m: &sync.RWMutex{},
			},
			args: args{
				nodeID: "1",
			},
			wantCount: 1,
		},
		{
			name: "remove one",
			fields: fields{
				UpstreamPool: []Node{
					&Upstream{nodeID: "1"},
					&Upstream{nodeID: "2"},
					&Upstream{nodeID: "3"},
				},
				m: &sync.RWMutex{},
			},
			args: args{
				nodeID: "2",
			},
			wantCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				UpstreamPool: tt.fields.UpstreamPool,
				load:         tt.fields.load,
				Policy:       tt.fields.Policy,
				m:            tt.fields.m,
			}
			if b.Remove(tt.args.nodeID); len(b.UpstreamPool) != tt.wantCount {
				t.Errorf("len(b.UpstreamPool) = %v, want5 %v", len(b.UpstreamPool), tt.wantCount)
			}
		})
	}
}
