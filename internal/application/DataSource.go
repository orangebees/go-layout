package application

import (
	"github.com/orangebees/go-oneutils/ExtErr"
)

type DataSource interface {
	GetById(id int) string
	Get() ExtErr.Err
}
