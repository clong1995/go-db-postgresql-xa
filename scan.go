package db

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/clong1995/go-ansi-color"
	"github.com/jackc/pgx/v5"
)

func scan[T any](rows pgx.Rows) (res []T, err error) {
	if !rows.Next() {
		res = make([]T, 0)
		return
	}

	var obj T
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Struct {
		if res, err = pgx.CollectRows[T](rows, pgx.RowToStructByPos[T]); err != nil {
			fmt.Println(pcolor.Err("CollectRows error: %v", err))
			return
		}
	} else {
		for rows.Next() {
			if err = rows.Scan(&obj); err != nil {
				log.Println(pcolor.Error(err))
				return
			}
			res = append(res, obj)
		}
	}

	return
}

func scanOne[T any](rows pgx.Rows) (res T, err error) {
	if res, err = pgx.CollectOneRow[T](
		rows,
		pgx.RowToStructByPos[T],
	); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Println(pcolor.Error(err))
			return
		}
	}

	return
}
