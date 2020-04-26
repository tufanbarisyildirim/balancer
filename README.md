# Balancer
Balancing machine that can be plugged anywhere

Balancer helps you to balance jobs/requests/messages between worker.  It is well tested, ready to be used in production under high throughput

```
-- 2,3 GHz 8-Core Intel Core i9

BenchmarkNextRoundRobin-16              26272492                44.2 ns/op
BenchmarkNextHash-16                     8979028               136 ns/op
BenchmarkNextLeastConnection-16         24265574                43.2 ns/op
BenchmarkNextLeastTime-16               10817994               112 ns/op
```