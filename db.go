package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/clong1995/go-ansi-color"
	"github.com/clong1995/go-config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pools map[string]*pgxpool.Pool
var mainPool *pgxpool.Pool

func init() {
	num, err := strconv.ParseInt(config.Value("MAXCONNS"), 10, 32)
	if err != nil {
		log.Fatalln(pcolor.Error(err))
		return
	}
	maxConn := int32(num)

	pools = make(map[string]*pgxpool.Pool)
	ds := config.Value("DATASOURCE")
	for _, v := range strings.Split(ds, ",") {
		var conf *pgxpool.Config
		if conf, err = pgxpool.ParseConfig(v); err != nil {
			log.Fatalln(pcolor.Error(err))
		}
		conf.MaxConns = maxConn

		var pool *pgxpool.Pool
		if pool, err = pgxpool.NewWithConfig(context.Background(), conf); err != nil {
			log.Fatalln(pcolor.Error(err))
		}

		if err = pool.Ping(context.Background()); err != nil {
			log.Fatalln(pcolor.Error(err))
		}

		if mainPool == nil {
			mainPool = pool
		}
		database := conf.ConnConfig.Database
		pools[database] = pool

		fmt.Println(pcolor.Succ("[PostgreSQL] conn %s", database))
	}
}

func Close() {
	for _, v := range pools {
		v.Close()
		fmt.Println(pcolor.Succ("[PostgreSQL] %s closed", v.Config().ConnConfig.Database))
	}
}
