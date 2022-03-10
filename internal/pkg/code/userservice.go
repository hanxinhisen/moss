// Created by Hisen at 2022/3/3.
package code

//go:generate codegen -type=int

const (
	// ErrUserNotFound - 404: User Not Found.
	ErrUserNotFound int = iota + 110001
	// ErrUserNotFound - 404: User Already Exist.
	ErrUserAlreadyExist
)
