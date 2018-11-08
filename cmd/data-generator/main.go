package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/linuxerwang/dgraph-bench/tasks"
)

const (
	maxUid           = 10000000
	maxDirectFriends = 1000

	k = 100
)

var (
	output = flag.String("output", "out.rdf.gz", "Output .gz file")
)

func main() {
	flag.Parse()
	f, err := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	w := gzip.NewWriter(f)
	for i := 1; i <= maxUid; i++ {
		meID := fmt.Sprintf("_:m.%d", i)
		writeNQuad(w, meID, "xid", fmt.Sprintf("\"%d\"", i))
		writeNQuad(w, meID, "name", fmt.Sprintf("\"%s\"", tasks.RandString(10)))
		writeNQuad(w, meID, "age", fmt.Sprintf("\"%d\"", 18+rand.Intn(80)))
		writeNQuad(w, meID, "created_at", fmt.Sprintf("\"%d\"", time.Now().UnixNano()))
		writeNQuad(w, meID, "updated_at", fmt.Sprintf("\"%d\"", time.Now().UnixNano()))

		friendCnt := randomNum()
		for j := 1; j <= friendCnt; j++ {
			fID := rand.Intn(maxUid)
			for fID == i {
				fID = rand.Intn(maxUid)
			}
			writeNQuad(w, meID, "friend_of", fmt.Sprintf("<_:m.%d>", fID))
		}
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func randomNum() int {
	// N(t) = N0 * e^(-k*t)
	return 5 + int(maxDirectFriends*math.Exp(-k*rand.Float64()))
}

func writeNQuad(w *gzip.Writer, s, p, o string) {
	str := fmt.Sprintf("<%v> <%v> %v .\n", s, p, o)
	if _, err := w.Write([]byte(str)); err != nil {
		panic(err)
	}
}
