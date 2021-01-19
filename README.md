# Balancer
Experimental, Zero allocation Balancing machine that can be plugged anywhere

Balancer helps you to balance jobs/requests/messages between workers.  It is well tested, ready to be used in production under high throughput. 

*Balancer is not a reverse proxy, it is a router, so you balance things between any type of Node Pool.


### API 
The interface of balancer is as simple as,

```go
//Balancer select a node to send load
type Balancer struct {
   Policy  SelectionPolicy 
   func Add(node ...Node) {}
   func Next(clientID string) Node {}
}
```

Selector is one of selection policies:

- [RoundRobin](roundrobin.go) - selects next available upstream on every request
- [Hash](hash.go) - matches client using consistent hashing by its own id (any string like ip address or user id)
- [LeastConnection](leastconnection.go) - selects the node that has the lowest active connection (using Node.Load)
- [LeastTime](leasttime.go) - selects the node that has the lowest response time (using Node.AverageResponseTime)

### Experiments
- [RoundRobinChannel](roundrobinchannel.go) - selects next available upstream using go channels as queue

Any type of object that satisfies Node Interface will work as node

```go
type Node interface {
	IsHealthy() bool
	TotalRequest() uint64
	AverageResponseTime() time.Duration
	Load() int64
	Host() string
}
```


You decide the selection policy on init like;
```go
package main 

balancer := Balancer{
    Policy:     &RoundRobin{},
}
```

Add nodes and balancer machine will decide the next upstream based on the policy

```go
balancer.Add(&Upstream{ Host:"worker-1" })
selectednode := balancer.Next("client-1")
```

### Benchmark
```
-- 2,3 GHz 8-Core Intel Core i9

➜ $ make bench
go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/tufanbarisyildirim/balancer
BenchmarkNextRoundRobin-16              19688180                52.7 ns/op             0 B/op          0 allocs/op
BenchmarkNextHash-16                    32723619                33.9 ns/op             0 B/op          0 allocs/op
BenchmarkNextLeastConnection-16         17938972                59.0 ns/op             0 B/op          0 allocs/op
BenchmarkNextLeastTime-16               11803462                98.3 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/tufanbarisyildirim/balancer  5.009s
```

Those results are for 10 upstreams where half of them are down but still in pool. (see [balancer_test.go](balancer_test.go) for details)
I tried to keep mocking as a real example as possible (Increasing current load, load times and kep half of them down all the time). So those results are worst case of finding the best upstream in pool.

### Real World Examples
-- coming soon

### TODO
- [ ] refer repos use balancer
- [ ] implement a real http load balancer as an example
- [ ] implement more policies

### [Contributing](CONTRIBUTING)
### [License](LICENSE)