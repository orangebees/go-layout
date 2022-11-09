package mysql

import "github.com/jmoiron/sqlx"

type Client struct {
	db          *sqlx.DB
	tablePrefix string
}

func New(db *sqlx.DB, tablePrefix string) *Client {
	return &Client{
		db:          db,
		tablePrefix: tablePrefix,
	}
}
