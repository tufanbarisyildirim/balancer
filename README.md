# Balancer
Balancing machine that can be plugged anywhere

Balancer helps you to balance jobs/requests/messages between workers.  It is well tested, ready to be used in production under high throughput. 

*Balancer is not a reverse proxy, it is a router so you balance things between any type of UpstreamPool.


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

- RoundRobin
- Hash
- LeastConnection
- LeastTime

Any type of object that satisfies Node Interface will work as node

```go
type Node interface {
	IsHealthy() bool
	GetTotalRequest() uint64
	GetAverageResponseTime() time.Duration
	GetLoad() int64
	GetHost() string
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
BenchmarkNextRoundRobin-16              25859118                44.1 ns/op
BenchmarkNextHash-16                     8691072               133 ns/op
BenchmarkNextLeastConnection-16         25841780                43.3 ns/op
BenchmarkNextLeastTime-16               12679124                88.6 ns/op
PASS
ok      github.com/tufanbarisyildirim/balancer  5.009s
```

Those results are for 10 upstreams where half of them are down but still in pool. (see balancer_test.go for details)
I tired to keep mocking as a real example as possible (Increasing current load, load times and kep half of them down all the time). So those results are worst case of finding the best upstream in pool.


### [Contributing](CONTRIBUTING)
### [License](LICENSE)