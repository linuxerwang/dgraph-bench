package tasks

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/dgo"
)

type BenchmarkCase func(dgraphCli *dgo.Dgraph) error

var (
	BenchTasks = map[string]BenchmarkCase{}
)

func ExecTask(name string, bc BenchmarkCase, dgraphCli *dgo.Dgraph, concurrency int) {
	count := int64(0)
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("!!!Error!!! %v\n", r)
						}
					}()

					start := time.Now()
					err := bc(dgraphCli)
					d := time.Since(start)

					status := "OK"
					if err != nil {
						status = "ERROR"
					} else {
						cnt := atomic.AddInt64(&count, 1)
						if cnt > 0 && cnt%1000 == 0 {
							fmt.Printf("Finished %d %s.\n", cnt, name)
						}
					}
					counters.WithLabelValues(name, status).Inc()
					durations.WithLabelValues(name, status).Observe(d.Seconds())
				}()
			}
		}()
	}
}
