package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

var (
	personXid int64
)

func init() {
	BenchTasks["insert-person"] = InsertPerson
}

func InsertPerson(dgraphCli *dgo.Dgraph) error {
	start := time.Now()

	xid := strconv.FormatInt(atomic.AddInt64(&personXid, 1), 10)

	person := &Person{
		Uid:       "_:" + xid,
		Xid:       xid,
		Name:      randString(10),
		CreatedAt: start.Unix(),
		UpdatedAt: start.Unix(),
	}

	// Mutate people node
	payload, err := json.Marshal(person)
	if err != nil {
		counters.WithLabelValues("insert-person", "ERROR").Inc()
		fmt.Printf("Failed to marshal person object: %v\n", err)
		return err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   payload,
	}

	txn := dgraphCli.NewTxn()
	as, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		counters.WithLabelValues("insert-person", "ERROR").Inc()
		fmt.Printf("Failed to call mutate: %v\n", err)
                return err
	}

	counters.WithLabelValues("insert-person", "OK").Inc()
	durations.WithLabelValues("insert-person", "OK").Observe(time.Since(start).Seconds())
	_ = as

	return nil
}
