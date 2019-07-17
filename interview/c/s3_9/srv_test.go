package srv

import (
	"io/ioutil"
	"net"
	"testing"
)

/*
go test -benchtime=10s -run=^$ -bench=^(BenchmarkNetworkRequest2)$

startNetworkDaemon()
goos: linux
goarch: amd64
BenchmarkNetworkRequest-8             10        1002300910 ns/op
PASS
ok      _/mnt/c/Users/admin/Desktop/go-interview/c/s3_9 11.042s

startNetworkDaemon2()
goos: linux
goarch: amd64
BenchmarkNetworkRequest-8           2000           5052589 ns/op
PASS
ok      _/mnt/c/Users/admin/Desktop/go-interview/c/s3_9 29.520s

*/

func BenchmarkNetworkRequest(b *testing.B) {
	bind := "localhost:8081"
	daemon := startNetworkDaemon(bind)
	daemon.Wait()

	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", bind)
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}

func BenchmarkNetworkRequest2(b *testing.B) {
	bind := "localhost:8080"
	daemon := startNetworkDaemon2(bind)
	daemon.Wait()

	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", bind)
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
