package service

import "github.com/orangebees/go-oneutils/ExtErr"

type DataSourceMock struct {
}

func (d DataSourceMock) Get() ExtErr.Err {
	e := ExtErr.NewErr("")
	return e
}

func (d DataSourceMock) NoneFunc() {
}

func (d DataSourceMock) GetById(id int) string {
	println(id)
	return ""
}

func NewMock() DataSourceMock {
	return DataSourceMock{}
}
