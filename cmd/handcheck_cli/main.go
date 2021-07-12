package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/nik0sc/mj"
	"github.com/nik0sc/mj/handcheck"
	"github.com/nik0sc/mj/wait"
)

func main() {
	in, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("scan: %s\n", err.Error())
		return
	}

	in = in[:len(in)-1]

	h, err := mj.ParseHand(in)
	if err != nil {
		fmt.Printf("cannot parse hand: %s\n", err.Error())
		return
	}

	sort.Sort(h)
	fmt.Printf("marshal: %x\n", h.Marshal())

	r := handcheck.OptHandRLEChecker{Split: false, UseMemo: true}.Check(h)
	fmt.Printf("solution: %s\n", r.String())

	waits := wait.Find(r, true)
	if len(waits) == 0 {
		fmt.Println("no waits")
	} else {
		fmt.Printf("waits: %s\n", mj.Hand(waits).String())
	}
}
