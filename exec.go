package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/clong1995/go-ansi-color"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

// QueryRow 查询一条
func QueryRow(db, query string, args ...any) (row pgx.Row, err error) {
	pool, err := singlePool(db)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	row = pool.QueryRow(context.Background(), query, args...)
	return
}

// QueryRowScan 查询并扫描
func QueryRowScan[T any](db, query string, args ...any) (res T, err error) {
	pool, err := singlePool(db)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	//QueryRow的源码也是调用了Query
	rows, err := pool.Query(context.Background(), query, args...)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	if res, err = scanOne[T](rows); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// Exec 执行
func Exec(db, query string, args ...any) (result pgconn.CommandTag, err error) {
	pool, err := singlePool(db)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	if result, err = pool.Exec(context.Background(), query, args...); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// Query 查询
func Query(db, query string, args ...any) (rows pgx.Rows, err error) {
	pool, err := singlePool(db)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	if rows, err = pool.Query(context.Background(), query, args...); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// QueryScan 查询并扫描
func QueryScan[T any](db, query string, args ...any) (res []T, err error) {
	rows, err := Query(db, query, args...)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	defer rows.Close()
	if res, err = scan[T](rows); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

func singlePool(db string) (pool *pgxpool.Pool, err error) {
	if db == "" {
		pool = mainPool
	} else {
		pool = pools[db]
		if pool == nil {
			err = errors.New(fmt.Sprintf("db[%s] is not exist", db))
			log.Println(pcolor.Error(err))
			return
		}
	}
	return
}
