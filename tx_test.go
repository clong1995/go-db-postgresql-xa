package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestBatchTx(t *testing.T) {
	defer Close()
	type args struct {
		query string
		data  [][]any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "batch tx",
			args: args{
				query: "INSERT INTO demo (id,name) VALUES ($1,$2);",
				data: [][]any{
					{4, "d"},
					{5, nil},
					{6, "f"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			/*Txs([]string{"account","access"}, func(txs []pgx.Tx) (err error) {
				txAccount := txs[0]
				txAccess := txs[1]
				return
			})*/

			if err := Tx(func(tx pgx.Tx) (err error) {
				if err = BatchTx(tx, tt.args.query, tt.args.data); (err != nil) != tt.wantErr {
					t.Errorf("BatchTx() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}); err != nil {
				t.Errorf("Tx() error = %v", err)
				return
			}
		})
	}
}

func TestCopyTx(t *testing.T) {
	defer Close()
	type args struct {
		tableName   string
		columnNames []string
		data        [][]any
	}
	tests := []struct {
		name             string
		args             args
		wantRowsAffected int64
		wantErr          bool
	}{
		{
			name: "batch tx",
			args: args{
				tableName:   "demo",
				columnNames: []string{"id", "name"},
				data: [][]any{
					{10, "j"},
					{11, nil},
					{12, "l"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Tx(func(tx pgx.Tx) (err error) {
				gotRowsAffected, err := CopyTx(tx, tt.args.tableName, tt.args.columnNames, tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("CopyTx() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("CopyTx() gotRowsAffected = %v", gotRowsAffected)
				return
			}); err != nil {
				t.Errorf("Tx() error = %v", err)
				return
			}
		})
	}
}

func TestTxExec(t *testing.T) {
	defer Close()
	type args struct {
		query string
		args  []any
	}
	tests := []struct {
		name       string
		args       args
		wantResult pgconn.CommandTag
		wantErr    bool
	}{
		{
			name: "tx exec",
			args: args{
				query: "INSERT INTO demo (id,name) VALUES ($1,$2);",
				args:  []any{11, "h"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Txs([]string{"account", "access"}, func(txs []pgx.Tx) (err error) {
				accountTx := txs[0]
				accessTx := txs[1]
				if _, err = TxExec(accountTx, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
					t.Errorf("TxExec() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				err = errors.New("new error")
				if err != nil {
					t.Errorf("TxExec() error = %v", err)
					return
				}
				if _, err = TxExec(accessTx, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
					t.Errorf("TxExec() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}); err != nil {
				t.Errorf("Tx() error = %v", err)
				return
			}
		})
	}
}

func TestTxQueryScan(t *testing.T) {
	defer Close()
	type args struct {
		query string
		args  []any
	}
	type testCase[T any] struct {
		name    string
		args    args
		wantRes []T
		wantErr bool
	}
	tests := []testCase[field]{
		{
			name: "tx query scan",
			args: args{
				query: "SELECT id,name FROM demo WHERE id = $1",
				args:  []any{17},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Tx(func(tx pgx.Tx) (err error) {
				if _, err = TxExec(tx, "INSERT INTO demo (id,name) VALUES ($1,$2);", 17, "q"); (err != nil) != tt.wantErr {
					t.Errorf("TxExec() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				gotRes, err := TxQueryScan[field](tx, tt.args.query, tt.args.args...)
				if (err != nil) != tt.wantErr {
					t.Errorf("TxQueryScan() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				t.Logf("QueryScan() gotRes = %#v,", gotRes)
				return
			}); err != nil {
				t.Errorf("Tx() error = %v", err)
				return
			}
		})
	}
}

func TestTxQueryRow(t *testing.T) {
	defer Close()
	type args struct {
		query string
		args  []any
	}
	tests := []struct {
		name    string
		args    args
		wantRow pgx.Row
	}{
		{
			name: "tx query row",
			args: args{
				query: "SELECT id,name FROM demo WHERE id = $1",
				args:  []any{17},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Tx(func(tx pgx.Tx) (err error) {
				gotRow := TxQueryRow(tx, tt.args.query, tt.args.args...)
				var id int64
				var name string
				if err = gotRow.Scan(&id, &name); err != nil {
					t.Errorf("TxQueryRow() error = %v", err)
					return
				}

				t.Logf("TxQueryRow() gotRes = %v,%v", id, name)
				return
			}); err != nil {
				t.Errorf("Tx() error = %v", err)
				return
			}
		})
	}
}

func TestTxQueryRowScan(t *testing.T) {
	defer Close()
	type args struct {
		query string
		args  []any
	}
	type testCase[T any] struct {
		name    string
		args    args
		wantRes T
		wantErr bool
	}
	tests := []testCase[field]{
		{
			name: "tx query row scan",
			args: args{
				query: "SELECT id,name FROM demo WHERE id = $1",
				args:  []any{17},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Tx(func(tx pgx.Tx) (err error) {
				gotRes, err := TxQueryRowScan[field](tx, tt.args.query, tt.args.args...)
				if (err != nil) != tt.wantErr {
					t.Errorf("TxQueryRowScan() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("TxQueryRow() gotRes = %#v", gotRes)
				return
			}); err != nil {
				t.Errorf("Tx() error = %v", err)
				return
			}
		})
	}
}
