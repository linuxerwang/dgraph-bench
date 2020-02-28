package tasks

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/dgo"
)

// BenchmarkCase ...
type BenchmarkCase func(dgraphCli *dgo.Dgraph, r *rand.Rand) error

var (
	BenchTasks = map[string]BenchmarkCase{}
)

func report(name string, count *int64) {
	prev := atomic.LoadInt64(count)
	timeCount := 0
	for range time.Tick(1 * time.Second) {
		timeCount++
		cnt := atomic.LoadInt64(count)
		throughput.WithLabelValues(name, "OK").Set(float64(cnt - prev))
		fmt.Printf("Time elapsed: %d, Taskname: %s, Speed: %d\n", timeCount, name, cnt-prev)
		prev = cnt
	}
}

func ExecTask(name string, bc BenchmarkCase, dgraphCli *dgo.Dgraph, concurrency int) {
	count := int64(0)
	go report(name, &count)
	for i := 0; i < concurrency; i++ {
		go func() {
			r := rand.New(rand.NewSource(time.Now().Unix()))
			for {
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("!!!Error!!! %v\n", r)
						}
					}()

					start := time.Now()
					err := bc(dgraphCli, r)
					d := time.Since(start)

					status := "OK"
					if err != nil {
						status = "ERROR"
					} else {
						atomic.AddInt64(&count, 1)
					}
					counters.WithLabelValues(name, status).Inc()
					durations.WithLabelValues(name, status).Observe(d.Seconds())
				}()
			}
		}()
	}
}
