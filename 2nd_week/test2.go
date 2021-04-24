package main

import (
	"database/sql"
	"errors"
	"fmt"

	pkg_errors "github.com/pkg/errors"
)

func doSthFailed(sqlText string) error {
	return fmt.Errorf("main: no data matched %v: %w", sqlText, sql.ErrNoRows)
}

func querySth(sqlText string) error {
	err := doSthFailed(sqlText)
	if err != nil {
		return pkg_errors.Wrap(err, "query failed")
	}
	return err
}

func task() error {
	// query sth
	sqlText := "SELECT * FROM alerts"
	err := querySth(sqlText)
	if errors.Is(err, sql.ErrNoRows) {
		// 没找到数据挺正常的，降级
		fmt.Printf("sql no rows: %T %v\n", pkg_errors.Cause(err), pkg_errors.Cause(err))
	} else if err != nil {
		fmt.Printf("unknown error! error info: %T %v\n stack trace:\n%+v\n",
			pkg_errors.Cause(err), pkg_errors.Cause(err), err)
		return err
	}
	return nil
}

func main() {
	task()
}
