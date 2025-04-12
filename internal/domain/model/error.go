package model

import "errors"

// ドメインエラー定義
var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound = errors.New("your requested item is not found")
	ErrConflict = errors.New("your item already exists")
	ErrBadRequest = errors.New("bad request")
	ErrInvalidCredentials = errors.New("invalid credentials")
)