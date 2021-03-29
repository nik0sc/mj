// +build js,wasm,go1.16

package main

import (
	"fmt"
	"syscall/js"

	"github.com/nik0sc/mj"
	"github.com/nik0sc/mj/handcheck"
)

//export optCheck
func optCheck(hand string, cb func(string, string), split bool, memo bool) {
	go func() {
		h, err := mj.ParseHand(hand)
		if err != nil {
			cb("", fmt.Sprintf("cannot parse hand: %s", err.Error()))
			return
		}

		r := handcheck.OptChecker{Split: split, UseMemo: memo}.Check(h)
		cb(r.String(), "")
	}()
}

//export optCheckSync
func optCheckSync(hand string, split bool, memo bool) string {
	h, err := mj.ParseHand(hand)
	if err != nil {
		return err.Error()
	}

	r := handcheck.OptChecker{Split: split, UseMemo: memo}.Check(h)
	return r.String()
}

type checker interface {
	Check(hand mj.Hand) mj.Group
}

func main() {
	var c chan struct{}
	js.Global().Set("checkHand", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_ = this

		alg := args[0].String()
		hand := args[1].String()
		cb := args[2]
		if cb.Type() != js.TypeFunction {
			panic(fmt.Sprintf("pos arg 1 is not Function, it is %s\n", cb.Type().String()))
		}
		split := args[3].Bool()
		memo := args[4].Bool()

		var c checker
		switch alg {
		case "opt":
			c = handcheck.OptChecker{Split: split, UseMemo: memo}
		case "optcnt":
			c = handcheck.OptCountChecker{Split: split, UseMemo: memo}
		case "greedy":
			c = handcheck.GreedyChecker{Split: split}
		default:
			panic("unrecognised alg: " + alg)
		}

		go func() {
			h, err := mj.ParseHand(hand)
			if err != nil {
				cb.Invoke(js.Null(), fmt.Sprintf("cannot parse hand: %s", err.Error()))
				return
			}

			cb.Invoke(c.Check(h).String(), js.Null())
		}()
		return nil
	}))
	<-c
}
