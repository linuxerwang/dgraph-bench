package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/linuxerwang/dgraph-bench/tasks"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

var (
	flagServers = flag.String("servers", "", "Comma separated dgraph server endpoints")

	dgraphCli *dgo.Dgraph
)

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

type BenchCase struct {
	Name        string `yaml:"name"`
	Concurrency int    `yaml:"concurrency"`
}

type Config struct {
	BenchCases []*BenchCase `yaml:"bench_cases"`
}

func loadConfig(fn string) *Config {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	conf := Config{}
	if err := yaml.Unmarshal([]byte(b), &conf); err != nil {
		log.Fatalf("error: %v", err)
	}
	return &conf
}

func main() {
	flag.Usage = usage
	flag.Parse()

	go tasks.StartPrometheusServer(3500)

	if *flagServers == "" {
		fmt.Println("Flag --servers is required.")
		os.Exit(2)
	}
	dgraphCli = connect(*flagServers)

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	conf := loadConfig(flag.Arg(0))
	fmt.Printf("Found %d tasks.\n", len(conf.BenchCases))

	for _, bench := range conf.BenchCases {
		task, ok := tasks.BenchTasks[bench.Name]
		if !ok {
			log.Fatalf("Task not found: %s\n", bench.Name)
		}
		fmt.Printf("Starting task %s (%d concurrent goroutines)\n", bench.Name, bench.Concurrency)
		go tasks.ExecTask(bench.Name, task, dgraphCli, bench.Concurrency)
	}

	select {}
}
