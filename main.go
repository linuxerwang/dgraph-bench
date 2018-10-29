package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

const (
	TypePerson = iota
)

var (
	flagServers     = flag.String("servers", "", "Comma separated dgraph server endpoints")
	flagConcurrency = flag.Int("c", 10, "concurrency")

	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	cnt = int64(0)

	dgraphCli *dgo.Dgraph
)

type Person struct {
	Name      string
	Xid       string
	Type      int
	CreatedAt int64
	UpdatedAt int64
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("\tdgraph-bench")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
}

func connect(servers string) *dgo.Dgraph {
	clis := make([]api.DgraphClient, 0, 5)
	for _, s := range strings.Split(strings.Replace(servers, " ", "", -1), ",") {
		if len(s) > 0 {
			fmt.Printf("Connect to server %s\n", s)
			conn, err := grpc.Dial(s, grpc.WithInsecure())
			if err != nil {
				panic(err)
			}
			clis = append(clis, api.NewDgraphClient(conn))
		}
	}
	return dgo.NewDgraphClient(clis...)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	go startPrometheusServer(3500)

	if *flagServers == "" {
		fmt.Println("Flag --servers is required.")
		os.Exit(2)
	}
	dgraphCli = connect(*flagServers)

	for i := 0; i < *flagConcurrency; i++ {
		go func() {
			for {
				start := time.Now()

				person := &Person{
					Xid:       strconv.FormatInt(rand.Int63n(1000)+1, 10),
					Name:      randString(10),
					CreatedAt: start.Unix(),
					UpdatedAt: start.Unix(),
				}

				// Mutate people node
				payload, err := json.Marshal(person)
				if err != nil {
					counters.WithLabelValues("insert-person", "ERROR").Inc()
					fmt.Printf("Failed to marshal person object: %v\n", err)
					return
				}

				mu := &api.Mutation{
					CommitNow: true,
					SetJson:   payload,
				}

				ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
				txn := dgraphCli.NewTxn()
				as, err := txn.Mutate(ctx, mu)
				if err != nil {
					counters.WithLabelValues("insert-person", "ERROR").Inc()
					fmt.Printf("Failed to call mutate: %v\n", err)
				}

				counters.WithLabelValues("insert-person", "OK").Inc()
				durations.WithLabelValues("insert-person", "OK").Observe(time.Since(start).Seconds())
				_ = as

				c := atomic.AddInt64(&cnt, 1)
				if c%1000 == 0 {
					fmt.Printf("Finished %d requests\n", c)
				}

				time.Sleep(20 * time.Millisecond)
			}
		}()
	}

	select {}
}
