package db

import (
	"log"
	"testing"

	"github.com/clong1995/go-ansi-color"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func TestQueryRow(t *testing.T) {
	defer Close()
	type args struct {
		query string
		args  []any
	}
	tests := []struct {
		name    string
		args    args
		wantRow pgx.Row
		wantErr bool
	}{
		{
			name: "QueryRow",
			args: args{
				query: "SELECT id,name FROM role WHERE id = $1",
				args: []any{
					160858030168911872,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRow, err := QueryRow(tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var id int64
			var name string
			if err = gotRow.Scan(&id, &name); err != nil {
				t.Error(err)
				return
			}
			t.Logf("QueryRow() gotRow = %v,%v", id, name)
		})
	}
}

func TestQueryRowScan(t *testing.T) {
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
			name: "QueryRowScan",
			args: args{
				query: "SELECT id,name FROM role WHERE id = $1",
				args: []any{
					160858030168911872,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := QueryRowScan[field](tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryRowScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("QueryRowScan() gotRes = %#v", gotRes)

		})
	}
}

func TestQuery(t *testing.T) {
	defer Close()
	type args struct {
		query string
		args  []any
	}
	tests := []struct {
		name     string
		args     args
		wantRows pgx.Rows
		wantErr  bool
	}{
		{
			name: "Query",
			args: args{
				query: "SELECT id,name FROM role",
				args:  []any{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, err := Query(tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			res := make([][]any, 0)
			for gotRows.Next() {
				var id int64
				var name string
				if err = gotRows.Scan(&id, &name); err != nil {
					log.Println(pcolor.Error(err))
					return
				}
				res = append(res, []any{id, name})
			}

			t.Logf("Query() gotRows = %v", res)
		})
	}
}

func TestQueryScan(t *testing.T) {
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
			name: "QueryScan",
			args: args{
				query: "SELECT id,name FROM role",
				args:  []any{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := QueryScan[field](tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("QueryScan() gotRes = %#v,", gotRes)
		})
	}
}

type field struct {
	Id   int64
	Name string
}

func TestExec(t *testing.T) {
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
			name: "Exec",
			args: args{
				query: "INSERT INTO role (id, name,description,key) VALUES ($1, $2,'测试','sc')",
				args:  []any{162040587868221441, "测试"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := Exec(tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Exec() gotResult = %v", gotResult)
		})
	}
}
