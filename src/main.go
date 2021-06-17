package main

import (
	"fibonacci/models"
	"fmt"
	logger "github.com/buhduh/go-logger"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	FIB_AT_I_REQ          string = "fibati"
	FIB_LESS_THAN_VAL_REQ        = "getmemoized"
	CLEAR_REQ                    = "clear"
)

var PORT string

var myLogger logger.Logger
var fibModel models.FibonacciModel

func calculateFibonacci(start models.FibonacciPair, index uint) []models.FibonacciValue {
	var curr *models.FibonacciValue = &start[1]
	toRet := []models.FibonacciValue{start[0], start[1]}
	for curr.Index < index {
		curr.Index++
		curr.Value = toRet[len(toRet)-2].Value + toRet[len(toRet)-1].Value
		toRet = append(toRet, *curr)
	}
	return toRet[2:]
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	myLogger.Tracef("request query: %c", vals)
	var val []string
	var ok bool
	if val, ok = vals["type"]; !ok || len(val) != 1 {
		http.Error(
			w, "the request query string requires 'type' to proceed and be singular",
			http.StatusBadRequest,
		)
		return
	}
	switch val[0] {
	case FIB_AT_I_REQ:
		if val, ok = vals["i"]; !ok {
			http.Error(
				w,
				fmt.Sprintf(
					"request type='%s' requires query parameter i to proceed", FIB_AT_I_REQ,
				),
				http.StatusBadRequest,
			)
			return
		}
		i, err := strconv.Atoi(val[0])
		if err != nil || i < 0 {
			http.Error(
				w,
				fmt.Sprintf(
					"request type='%s' requires query parameter i to be a positive integer", FIB_AT_I_REQ,
				),
				http.StatusBadRequest,
			)
		}
		fibVals, err := fibModel.GetValues(uint(i))
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf(
					"failed retrieving memoized results with error: '%s'", err,
				),
				http.StatusInternalServerError,
			)
			return
		}
		//probably dumb, really inefficient
		fibVals = append([]models.FibonacciValue{models.INITIAL_FIBONACCI_PAIR[1]}, fibVals...)
		fibVals = append([]models.FibonacciValue{models.INITIAL_FIBONACCI_PAIR[0]}, fibVals...)
		fib := calculateFibonacci(
			models.FibonacciPair{
				fibVals[len(fibVals)-2],
				fibVals[len(fibVals)-1],
			}, uint(i),
		)
		myLogger.Infof("fibVals len: %d", len(fibVals))
		err = fibModel.PutValues(fib)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf(
					"failed inserting values into db with error: '%s'", err,
				),
				http.StatusInternalServerError,
			)
			return
		}
		var toWrite models.FibonacciValue
		if len(fib) == 0 {
			toWrite = fibVals[len(fibVals)-1]
		} else {
			toWrite = fib[len(fib)-1]
		}
		io.WriteString(w, toWrite.String())
		return
	case FIB_LESS_THAN_VAL_REQ:
		if val, ok = vals["value"]; !ok || len(val) != 1 {
			http.Error(
				w,
				fmt.Sprintf(
					"request type='%s' requires query parameter value and be singular to proceed",
					FIB_LESS_THAN_VAL_REQ,
				),
				http.StatusBadRequest,
			)
			return
		}
		iVal, err := strconv.Atoi(val[0])
		if err != nil || iVal < 0 {
			http.Error(
				w,
				fmt.Sprintf(
					"could not convert value parameter to a positive integer with error: '%s'", err,
				),
				http.StatusBadRequest,
			)
			return
		}
		numLessThan, err := fibModel.GetValuesLessThan(uint(iVal))
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf(
					"failed retrieving memoized values less than: %d, with error: '%s'", iVal, err,
				),
				http.StatusInternalServerError,
			)
			return
		}
		//Don't forget, i don't save first two in db
		if *numLessThan > 0 {
			*numLessThan += uint(2)
		}
		io.WriteString(
			w,
			fmt.Sprintf(
				"there are %d memoized values less than %d", *numLessThan, iVal,
			),
		)
		return
	case CLEAR_REQ:
		err := fibModel.Clear()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf(
					"failed clearing database with error: '%s'", err,
				),
				http.StatusInternalServerError,
			)
			return
		}
		io.WriteString(w, "database cleared")
		return
	}
}

func main() {
	myLogger = logger.NewLogger(logger.TRACE, "main logger")
	fibModel = models.NewFibonacciModel(myLogger)
	http.HandleFunc("/fib", mainHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil); err != nil {
		myLogger.Fatalf("http server failed with errror: '%s'", err)
		os.Exit(1)
	}
}
