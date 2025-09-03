package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/clong1995/go-ansi-color"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*// QueryRow 查询一条
func QueryRow(query string, args ...any) (row pgx.Row, err error) {
	return QueryRowDB("", query, args...)
}

func QueryRowDB(db, query string, args ...any) (row pgx.Row, err error) {
	pool, err := singlePool(db)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	row = pool.QueryRow(context.Background(), query, args...)
	return
}

// QueryRowScan 查询并扫描
func QueryRowScan[T any](query string, args ...any) (res T, err error) {
	return QueryRowScanDB[T]("", query, args...)
}

func QueryRowScanDB[T any](db, query string, args ...any) (res T, err error) {
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
	//CollectOneRow执行了close
	//defer rows.Close()
	if res, err = scanOne[T](rows); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}*/

// Exec 执行
func Exec(query string, args ...any) (result pgconn.CommandTag, err error) {
	return ExecDB("", query, args...)
}
func ExecDB(db, query string, args ...any) (result pgconn.CommandTag, err error) {
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
func Query(query string, args ...any) (rows pgx.Rows, err error) {
	return QueryDB("", query, args...)
}
func QueryDB(db, query string, args ...any) (rows pgx.Rows, err error) {
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
func QueryScan[T any](query string, args ...any) (res []T, err error) {
	return QueryScanDB[T]("", query, args...)
}
func QueryScanDB[T any](db, query string, args ...any) (res []T, err error) {
	rows, err := QueryDB(db, query, args...)
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
		pool = GetPool(db)
		if pool == nil {
			err = errors.New(fmt.Sprintf("db[%s] is not exist", db))
			log.Println(pcolor.Error(err))
			return
		}
	}
	return
}

func GetPool(db string) (pool *pgxpool.Pool) {
	pool = pools[db]
	return
}
