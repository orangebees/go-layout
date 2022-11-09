package service

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/orangebees/go-layout/internal/dao/mysql"
	"github.com/orangebees/go-layout/internal/dao/redis"
)

type DataSourceDefault struct {
	rdb *redis.Client
	db  *mysql.Client
}

func (d DataSourceDefault) NoneFunc() {
}

func (d DataSourceDefault) GetById(id int) string {
	return ""
}

func NewDefault(db *sqlx.DB, tablePrefix string, rdb *goredis.Client, keyPrefix string) (d DataSourceDefault) {
	d.db = mysql.New(db, tablePrefix)
	d.rdb = redis.New(rdb, keyPrefix)
	return
}
