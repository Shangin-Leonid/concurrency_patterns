package main

import (
	"context"
	"fmt"
)

func run_bridge_channel() {

	in1 := constIntGenerator(3, 111)
	in2 := constIntGenerator(3, 88)
	in3 := constIntGenerator(3, 2)

	chOfChs := make(chan (<-chan int))
	go func() {
		defer close(chOfChs)

		chOfChs <- in1
		chOfChs <- in2
		chOfChs <- in3
	}()

	concatCh := BridgeChannel(context.Background(), chOfChs)

	for v := range concatCh {
		fmt.Print(v, " ")
	}
	fmt.Println()

}

// BridgeChannel
//
// TODO:
//   - docs
func BridgeChannel[T any](ctx context.Context, chOfChs <-chan <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		var curIn <-chan T
		var ok bool

		resendFromCurIn := func() (needReturn bool) {
			var value T
			var ok bool

			for {
				select {
				case <-ctx.Done():
					return true
				case value, ok = <-curIn:
					if !ok {
						return false
					}
				}

				select {
				case <-ctx.Done():
					return true
				case out <- value:
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case curIn, ok = <-chOfChs:
				if !ok {
					return
				}
			}

			needReturn := resendFromCurIn()
			if needReturn {
				return
			}
		}
	}()

	return out
}
