package db

import (
	"testing"

	"github.com/jackc/pgx/v5"
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
