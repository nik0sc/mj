// +build js,wasm,go1.16
// +build !tinygo

package main

import (
	"fmt"
	"syscall/js"

	"mj"
	"mj/handcheck"
)

// window.optCheck(String, Function(String, String), Boolean, Boolean)
var optCheck js.Func
var stop chan struct{}

func init() {
	stop = make(chan struct{})
	optCheck = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		hand := args[0].String()
		cb := args[1]
		split := args[2].Bool()
		useMemo := args[3].Bool()

		if cb.Type() != js.TypeFunction {
			panic("pos arg 1 not function")
		}

		// This might take a while especially on mobile devices
		go func() {
			defer func() {
				re := recover()
				if re != nil {
					fmt.Printf("panic!: %v\n", re)
					cb.Invoke(js.Null(), re.(error).Error())
				}
			}()

			h, err := mj.ParseHand(hand)
			if err != nil {
				cb.Invoke(js.Null(), err.Error())
				return
			}
			r := handcheck.OptChecker{Split: split, UseMemo: useMemo}.Check(h)
			cb.Invoke(r.String(), js.Null())
		}()
		return nil
	})
}

func main() {
	js.Global().Get("window").Set("optCheck", optCheck)
	<-stop
	js.Global().Delete("optCheck")
	optCheck.Release()
}
