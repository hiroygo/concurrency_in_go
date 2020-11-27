package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

func connectToExternalServer() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

func externalServerConnPool() *sync.Pool {
	p := &sync.Pool{
		New: connectToExternalServer,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New())
	}
	return p
}

func startFastServer() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		connPool := externalServerConnPool()

		server, err := net.Listen("tcp", "localhost:9090")
		if err != nil {
			panic(err)
		}
		defer server.Close()
		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				panic(err)
			}

			svConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svConn)
			// conn.Close() を defer にすると for 内なので切断されない
			conn.Close()
		}
	}()

	return &wg
}

func startSlowServer() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			panic(err)
		}
		defer server.Close()
		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				panic(err)
			}

			connectToExternalServer()
			fmt.Fprintln(conn, "")
			// conn.Close() を defer にすると for 内なので切断されない
			conn.Close()
		}
	}()

	return &wg
}

func init() {
	slowSvStarted := startSlowServer()
	slowSvStarted.Wait()

	fastSvStarted := startFastServer()
	fastSvStarted.Wait()
}

func BenchmarkRequestToFastServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:9090")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}

func BenchmarkRequestToSlowServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
