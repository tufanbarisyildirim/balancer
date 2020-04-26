package balancer

import (
	"testing"
	"time"
)

func TestUpstream_GetAverageResponseTime(t *testing.T) {
	type fields struct {
		Node             Node
		Weight           int
		Load             int64
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
				Load:             0,
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
				Load:             0,
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
			u := &Upstream{
				Node:             tt.fields.Node,
				Weight:           tt.fields.Weight,
				Load:             tt.fields.Load,
				RequestCount:     tt.fields.RequestCount,
				TotalRequestTime: tt.fields.TotalRequestTime,
				Host:             tt.fields.Host,
				Healthy:          tt.fields.Healthy,
			}
			if got := u.GetAverageResponseTime(); got != tt.want {
				t.Errorf("Upstream.GetAverageResponseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
