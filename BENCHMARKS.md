```bash
go test -bench=. -benchmem `go list ./... | grep -v vendor`
?       github.com/vigo5190/go-socks5   [no test files]
goos: darwin
goarch: amd64
pkg: github.com/vigo5190/go-socks5/proxy
BenchmarkProxy_Start-8             10000            161408 ns/op            5131 B/op         66 allocs/op
BenchmarkAuth_Auth-8             3000000               441 ns/op             104 B/op          4 allocs/op
PASS
```