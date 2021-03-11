// +build js,wasm,go1.16

package main

import (
	"fmt"
	"syscall/js"

	"mj"
	"mj/handcheck"
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

func main() {
	var c chan struct{}
	js.Global().Set("optCheck", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_ = this

		hand := args[0].String()
		cb := args[1]
		if cb.Type() != js.TypeFunction {
			panic(fmt.Sprintf("pos arg 1 is not Function, it is %s\n", cb.Type().String()))
		}
		split := args[2].Bool()
		memo := args[3].Bool()

		go func() {
			h, err := mj.ParseHand(hand)
			if err != nil {
				cb.Invoke(js.Null(), fmt.Sprintf("cannot parse hand: %s", err.Error()))
				return
			}

			r := handcheck.OptChecker{Split: split, UseMemo: memo}.Check(h)
			cb.Invoke(r.String(), js.Null())
		}()
		return nil
	}))

	js.Global().Set("optCountCheck", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_ = this

		hand := args[0].String()
		cb := args[1]
		if cb.Type() != js.TypeFunction {
			panic(fmt.Sprintf("pos arg 1 is not Function, it is %s\n", cb.Type().String()))
		}
		split := args[2].Bool()
		memo := args[3].Bool()

		go func() {
			h, err := mj.ParseHand(hand)
			if err != nil {
				cb.Invoke(js.Null(), fmt.Sprintf("cannot parse hand: %s", err.Error()))
				return
			}

			r := handcheck.OptCountChecker{Split: split, UseMemo: memo}.Check(h)
			cb.Invoke(r.String(), js.Null())
		}()
		return nil
	}))

	js.Global().Set("greedyCheck", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		_ = this

		hand := args[0].String()
		cb := args[1]
		if cb.Type() != js.TypeFunction {
			panic(fmt.Sprintf("pos arg 1 is not Function, it is %s\n", cb.Type().String()))
		}
		split := args[2].Bool()
		failfast := args[3].Bool()

		go func() {
			h, err := mj.ParseHand(hand)
			if err != nil {
				cb.Invoke(js.Null(), fmt.Sprintf("cannot parse hand: %s", err.Error()))
				return
			}

			r := handcheck.GreedyChecker{Split: split, FailFast: failfast}.Check(h)
			rstr := fmt.Sprintf("%+v", r)
			cb.Invoke(rstr, js.Null())
		}()
		return nil
	}))
	<-c
}
