package main

import (
	"fibonacci/models"
	"fmt"
	"testing"
)

var fibTests = []struct {
	start    models.FibonacciPair
	index    uint
	expStart models.FibonacciValue
	expEnd   models.FibonacciValue
	expLen   int
}{
	{
		models.INITIAL_FIBONACCI_PAIR,
		12,
		models.FibonacciValue{
			Index: 2,
			Value: 1,
		},
		models.FibonacciValue{
			Index: 12,
			Value: 144,
		},
		11,
	},
	{
		models.FibonacciPair{
			models.FibonacciValue{
				Index: 5,
				Value: 5,
			},
			models.FibonacciValue{
				Index: 6,
				Value: 8,
			},
		},
		15,
		models.FibonacciValue{
			Index: 7,
			Value: 13,
		},
		models.FibonacciValue{
			Index: 15,
			Value: 610,
		},
		9,
	},
}

func TestCalulateFibonacci(t *testing.T) {
	pMsg := func(i int, fStr string, data ...interface{}) string {
		return fmt.Sprintf(
			"calculate fibonacci, index: %d, %s",
			i, fmt.Sprintf(fStr, data...),
		)
	}
	for i, tCase := range fibTests {
		res := calculateFibonacci(tCase.start, tCase.index)
		if !res[0].Equals(tCase.expStart) {
			t.Errorf(
				pMsg(
					i, "calculateFibonacci returned list start not expected; expected %s, got %s",
					tCase.expStart, res[0],
				),
			)
		}
		if !res[len(res)-1].Equals(tCase.expEnd) {
			t.Errorf(
				pMsg(
					i, "calculateFibonacci returned list end not expected; expected %s, got %s",
					tCase.expEnd, res[len(res)-1],
				),
			)
		}
		if len(res) != tCase.expLen {
			t.Errorf(
				pMsg(
					i, "calculateFibonacci returned list length not expected; expected %d, got %d",
					tCase.expLen, len(res),
				),
			)

		}
	}
}
