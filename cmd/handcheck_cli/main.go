package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"mj"
	"mj/handcheck"
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

	r := handcheck.OptChecker{Split: false, UseMemo: true}.Check(h)
	fmt.Printf("solution: %s\n", r.String())
}
