package main

import (
	"fmt"
	"time"
)

func workWithHeartbeat(done <-chan interface{}, beatInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	beatCh := make(chan interface{})
	retCh := make(chan time.Time)

	go func() {
		defer close(beatCh)
		defer close(retCh)

		// Ticker は終わったら Stop() で開放すること
		beatTicker := time.NewTicker(beatInterval)
		workTicker := time.NewTicker(3 * beatInterval)
		defer beatTicker.Stop()
		defer workTicker.Stop()

		heartbeat := func() {
			select {
			case beatCh <- struct{}{}:
			// beatCh が受信されていない場合でも
			// デッドロックしないように default で抜けるようにしておく
			default:
			}
		}

		work := func(t time.Time) {
			for {
				select {
				case <-done:
					return
				case <-beatTicker.C:
					// done と同じように heartbeat も常に case に入れる

					// 本来の目的である仕事はまだできていないので
					// return で抜けない
					heartbeat()
				case retCh <- t:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-beatTicker.C:
				// done と同じように heartbeat も常に case に入れる
				heartbeat()
			case t := <-workTicker.C:
				work(t)
			}
		}
	}()

	return beatCh, retCh
}

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	beatCh, workCh := workWithHeartbeat(done, timeout/2)
	for {
		select {
		case _, ok := <-beatCh:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case w, ok := <-workCh:
			if !ok {
				return
			}
			fmt.Printf("work: %v\n", w)
		case <-time.After(timeout):
			fmt.Println("timeout")
			return
		}
	}
}
