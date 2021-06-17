package models

import (
	"database/sql"
	"fmt"
	logger "github.com/buhduh/go-logger"
	"strings"
)

type FibonacciValue struct {
	Index uint
	Value uint
}

type FibonacciPair [2]FibonacciValue

var INITIAL_FIBONACCI_PAIR FibonacciPair = FibonacciPair{
	FibonacciValue{
		Index: 0,
		Value: 0,
	},
	FibonacciValue{
		Index: 1,
		Value: 1,
	},
}

func (f FibonacciValue) String() string {
	return fmt.Sprintf("Fib(%d) = %d", f.Index, f.Value)
}

func (f1 FibonacciValue) Equals(f2 FibonacciValue) bool {
	return f1.Value == f2.Value && f1.Index == f2.Index
}

type FibonacciModel interface {
	GetValues(uint) ([]FibonacciValue, error)
	//assumed ordered
	PutValues([]FibonacciValue) error
	GetValuesLessThan(uint) (*uint, error)
	Clear() error
}

type fibModel struct {
	fibLogger logger.Logger
}

var clearStatement *sql.Stmt = nil

func (f *fibModel) Clear() error {
	var err error
	if err = safelyConnect(); err != nil {
		return err
	}
	if clearStatement == nil {
		clearStatement, err = connection.Prepare(
			"DELETE FROM fibonacci",
		)
		if err != nil {
			return err
		}
	}
	_, err = clearStatement.Exec()
	return err
}

var getValueLessStatement *sql.Stmt = nil

func (f *fibModel) GetValuesLessThan(val uint) (*uint, error) {
	var err error
	if err = safelyConnect(); err != nil {
		return nil, err
	}
	if getValueLessStatement == nil {
		getValueLessStatement, err = connection.Prepare(
			"SELECT COUNT (index) FROM fibonacci WHERE value <= $1",
		)
		if err != nil {
			return nil, err
		}
	}
	var rows *sql.Rows
	rows, err = getValueLessStatement.Query(val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var toRet uint
	rows.Next()
	err = rows.Scan(&toRet)
	if err != nil {
		return nil, err
	}
	f.fibLogger.Tracef("got less than: %d", toRet)
	if rows.NextResultSet() {
		return nil, fmt.Errorf("should only retrieve a single result from this query")
	}
	return &toRet, nil
}

func (f *fibModel) PutValues(vals []FibonacciValue) error {
	f.fibLogger.Infof("attempting to insert %d items", len(vals))
	if len(vals) == 0 {
		return nil
	}
	toPrep := make([]string, len(vals))
	for i := 0; i < len(toPrep)*2; i += 2 {
		toPrep[i/2] = fmt.Sprintf("($%d, $%d)", i+1, i+2)
	}
	//yes i know this is bad practice, not bothering for the sake of time
	prepStr := fmt.Sprintf(
		"INSERT INTO fibonacci (index, value) VALUES %s ON CONFLICT DO NOTHING",
		strings.Join(toPrep, ","),
	)
	stmt, err := connection.Prepare(prepStr)
	if err != nil {
		return err
	}
	toInsert := make([]interface{}, len(vals)*2)
	for i := 0; i < len(vals)*2; i += 2 {
		toInsert[i] = vals[i/2].Index
		toInsert[i+1] = vals[i/2].Value
		f.fibLogger.Infof("want to insert val: %s", vals[i/2])
	}
	_, err = stmt.Exec(toInsert...)
	return err
}

var getValueStatement *sql.Stmt = nil

func (f *fibModel) GetValues(i uint) ([]FibonacciValue, error) {
	var err error
	if err = safelyConnect(); err != nil {
		return nil, err
	}
	if getValueStatement == nil {
		getValueStatement, err = connection.Prepare(
			"SELECT index, value FROM fibonacci WHERE index <= $1 ORDER BY index ASC",
		)
		if err != nil {
			return nil, err
		}
	}
	var rows *sql.Rows
	rows, err = getValueStatement.Query(i)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var index, value uint
	toRet := make([]FibonacciValue, 0)
	for rows.Next() {
		f.fibLogger.Trace("looping through rows")
		err = rows.Scan(&index, &value)
		if err != nil {
			return nil, err
		}
		fibRet := FibonacciValue{
			Index: index,
			Value: value,
		}
		f.fibLogger.Infof("appending %s to return list in GetValues", fibRet)
		toRet = append(toRet, fibRet)
	}
	return toRet, nil
}

func NewFibonacciModel(myLogger logger.Logger) FibonacciModel {
	return &fibModel{
		fibLogger: myLogger,
	}
}
