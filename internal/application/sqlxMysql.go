package application

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

var sqlxMysqlClient *sqlx.DB
var tablePrefix = "app"

func GetSqlxMysqlClient() *sqlx.DB {
	if sqlxMysqlClient != nil {
		return sqlxMysqlClient
	}
	return createSqlxMysqlClient()
}
func GetSqlxMysqlTablePrefix() string {
	return tablePrefix
}
func createSqlxMysqlClient() *sqlx.DB {
	if v := os.Getenv("SQLX_MYSQL_TABLEPREFIX"); v != "" {
		tablePrefix = v
	}
	var sc sqlConfig
	sc.Host = "127.0.0.1"
	sc.Port = "3306"
	sc.UserName = "root"
	if v := os.Getenv("SQLX_MYSQL_HOST"); v != "" {
		sc.Host = v
	}
	if v := os.Getenv("SQLX_MYSQL_PORT"); v != "" {
		sc.Port = v
	}
	if v := os.Getenv("SQLX_MYSQL_USERNAME"); v != "" {
		sc.UserName = v
	}
	if v := os.Getenv("SQLX_MYSQL_PASSWORD"); v != "" {
		sc.Password = v
	}
	if v := os.Getenv("SQLX_MYSQL_DBNAME"); v != "" {
		sc.DbName = v
	}
	db, err := sqlx.Connect("mysql", sc.UserName+":"+sc.Password+"@tcp("+sc.Host+":"+sc.Port+")/"+sc.DbName+"?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		panic(err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		return nil
	}
	GetLogger().Info().Msg("mysql connect success")
	//注册释放方法
	closes = append(closes, db.Close)
	return db
}
