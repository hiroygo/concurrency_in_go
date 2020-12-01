package main

import (
	"fmt"
	"net/http"
)

type headResult struct {
	resp *http.Response
	err  error
}

func httpHeads(urls []string) <-chan headResult {
	c := make(chan headResult)
	go func() {
		defer close(c)
		for _, url := range urls {
			r, err := http.Head(url)
			c <- headResult{resp: r, err: err}
		}
	}()
	return c
}

func main() {
	heads := httpHeads([]string{"https://www.google.co.jp/", "badhost"})
	for h := range heads {
		if h.err != nil {
			fmt.Println(h.err)
			continue
		}
		fmt.Println(h.resp.Status)
	}
}
