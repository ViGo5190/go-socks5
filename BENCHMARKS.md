```bash
go test -bench=. -benchmem `go list ./... | grep -v vendor`
?       github.com/vigo5190/go-socks5   [no test files]
goos: darwin
goarch: amd64
pkg: github.com/vigo5190/go-socks5/proxy
BenchmarkProxy_Start-8             10000            145070 ns/op            5130 B/op         66 allocs/op
BenchmarkAuth_Auth-8             3000000               499 ns/op              96 B/op          3 allocs/op
PASS
```