// Created by Hisen at 2022/3/3.
package code

import (
	"github.com/marmotedu/errors"
	"github.com/novalagung/gubrak"
)

type ErrCode struct {
	C    int
	HTTP int
	Ext  string
	Ref  string
}

func (code ErrCode) HTTPStatus() int {
	return code.HTTP
}

func (code ErrCode) String() string {
	return code.Ext
}

func (code ErrCode) Reference() string {
	return code.Ref

}

var _ errors.Coder = &ErrCode{}

func (code ErrCode) Code() int {
	return code.C
}

func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	errors.MustRegister(coder)
}
