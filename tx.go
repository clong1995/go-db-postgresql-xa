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

// BatchTx 批量数据的插入
func BatchTx(tx pgx.Tx, query string, data [][]any) (err error) {
	batch := &pgx.Batch{}
	for _, v := range data {
		_ = batch.Queue(query, v...)
	}
	br := tx.SendBatch(context.Background(), batch)
	if err = br.Close(); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// CopyTx 超大量数据插入的
func CopyTx(tx pgx.Tx, tableName string, columnNames []string, data [][]any) (rowsAffected int64, err error) {
	table := pgx.Identifier{tableName}
	if rowsAffected, err = tx.CopyFrom(
		context.Background(),
		table,
		columnNames,
		pgx.CopyFromRows(data),
	); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// Tx 事物
func Tx(handle func(tx pgx.Tx) (err error)) (err error) {
	tx, err := mainPool.Begin(context.Background())
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				log.Println(pcolor.Error(rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(context.Background()); commitErr != nil {
				log.Println(pcolor.Error(commitErr))
			}
		}
	}()

	if err = handle(tx); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// Txs 事物
func Txs(dbs []string, handle func(txs []pgx.Tx) (err error)) (err error) {
	var ps []*pgxpool.Pool
	if dbs == nil || len(dbs) == 0 {
		ps = []*pgxpool.Pool{mainPool}
	} else {
		ps = make([]*pgxpool.Pool, len(dbs))
		for i, v := range dbs {
			p := pools[v]
			if p == nil {
				err = errors.New(fmt.Sprintf("db[%s] is not exist", v))
				log.Println(pcolor.Error(err))
				return
			}
			ps[i] = p
		}
	}

	txs := make([]pgx.Tx, len(dbs))
	for i, pool := range ps {
		//开启事物
		if txs[i], err = pool.Begin(context.Background()); err != nil {
			log.Println(pcolor.Error(err))
			return
		}
	}

	defer func() {
		if err != nil {
			for _, tx := range txs {
				if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
					log.Println(pcolor.Error(rollbackErr))
				}
			}
		} else {
			for _, tx := range txs {
				if commitErr := tx.Commit(context.Background()); commitErr != nil {
					log.Println(pcolor.Error(commitErr))
				}
			}
		}
	}()

	if err = handle(txs); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// TxExec 事物内执行
func TxExec(tx pgx.Tx, query string, args ...any) (result pgconn.CommandTag, err error) {
	if result, err = tx.Exec(context.Background(), query, args...); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

// TxQueryScan 事物内查询并扫描
func TxQueryScan[T any](tx pgx.Tx, query string, args ...any) (res []T, err error) {
	rows, err := TxQuery(tx, query, args...)
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

// TxQuery 事物内查询
func TxQuery(tx pgx.Tx, query string, args ...any) (rows pgx.Rows, err error) {
	if rows, err = tx.Query(context.Background(), query, args...); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}

/*// TxQueryRow 事物内查询
func TxQueryRow(tx pgx.Tx, query string, args ...any) (row pgx.Row) {
	return tx.QueryRow(context.Background(), query, args...)
}

func TxQueryRowScan[T any](tx pgx.Tx, query string, args ...any) (res T, err error) {
	rows, err := tx.Query(context.Background(), query, args...)
	if err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	if res, err = scanOne[T](rows); err != nil {
		log.Println(pcolor.Error(err))
		return
	}
	return
}*/
