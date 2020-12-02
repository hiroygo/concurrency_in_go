package main

import (
	"fmt"
	"net/http"
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
				c <- headResponse{resp: resp, err: err}
			}
		}
	}()
	return c
}

func main() {
	done := make(chan interface{})
	defer close(done)

	for r := range head(done, urlGenerator(done,
		"https://golang.org/pkg/net/http",
        "https://www.youtube.com/",
        "https://github.com/",
		"https://www.google.co.jp/")) {
		if r.err != nil {
			fmt.Println(r.err)
			continue
		}
		fmt.Println(r.resp.Status)
	}
}
