package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
)

func urlGenerator(done <-chan interface{}, urls ...string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for _, url := range urls {
			select {
			case <-done:
				return
			case c <- url:
			}
		}
	}()
	return c
}

type headResponse struct {
	resp *http.Response
	err  error
}

func head(done <-chan interface{}, urlCh <-chan string) <-chan headResponse {
	c := make(chan headResponse)
	go func() {
		defer close(c)
		for {
			select {
			case <-done:
				return
			case url, ok := <-urlCh:
				if !ok {
					return
				}
				resp, err := http.Head(url)

				select {
				case <-done:
					return
				case c <- headResponse{resp: resp, err: err}:
				}
			}
		}
	}()
	return c
}

func fanout(done <-chan interface{}, urlCh <-chan string) []<-chan headResponse {
	chans := make([]<-chan headResponse, runtime.NumCPU())
	for i := 0; i < len(chans); i++ {
		chans[i] = head(done, urlCh)
	}
	return chans
}

func fanin(done <-chan interface{}, chans ...<-chan headResponse) <-chan headResponse {
	var wg sync.WaitGroup
	multiplexedCh := make(chan headResponse)

	multiplex := func(in <-chan headResponse) {
		defer wg.Done()
		for resp := range in {
			select {
			case <-done:
				return
			case multiplexedCh <- resp:
			}
		}
	}

	wg.Add(len(chans))
	for _, c := range chans {
		go multiplex(c)
	}

	go func() {
		defer close(multiplexedCh)
		wg.Wait()
	}()

	return multiplexedCh
}

func main() {
	done := make(chan interface{})
	defer close(done)

	urlCh := urlGenerator(done,
		"https://golang.org/pkg/net/http",
		"https://www.youtube.com/",
		"https://github.com/",
		"https://www.google.co.jp/")
	chans := fanout(done, urlCh)
	respCh := fanin(done, chans...)
	for r := range respCh {
		if r.err != nil {
			fmt.Println(r.err)
			continue
		}
		fmt.Println(r.resp.Status)
	}
}
