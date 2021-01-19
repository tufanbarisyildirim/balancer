package balancer

import (
	"testing"
	"time"
)

func TestMockNode_AverageResponseTime(t *testing.T) {
	type fields struct {
		Node             Node
		Weight           int
		load             int64
		RequestCount     uint64
		TotalRequestTime uint64
		Host             string
		Healthy          bool
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "basic time calc",
			fields: fields{
				Node:             nil,
				Weight:           1,
				load:             0,
				RequestCount:     100,
				TotalRequestTime: 100 * 100,
				Healthy:          true,
				Host:             "127.0.0.1",
			},
			want: time.Duration(100 * time.Nanosecond),
		},
		{
			name: "basic time calc again :)",
			fields: fields{
				Node:             nil,
				Weight:           1,
				load:             0,
				RequestCount:     100,
				TotalRequestTime: uint64(100 * time.Second),
				Healthy:          true,
				Host:             "127.0.0.1",
			},
			want: time.Duration(time.Second),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &MockNode{
				Node:             tt.fields.Node,
				Weight:           tt.fields.Weight,
				load:             tt.fields.load,
				requestCount:     tt.fields.RequestCount,
				totalRequestTime: tt.fields.TotalRequestTime,
				nodeID:           tt.fields.Host,
				healthy:          tt.fields.Healthy,
			}
			if got := u.AverageResponseTime(); got != tt.want {
				t.Errorf("MockNode.AverageResponseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
