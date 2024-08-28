package router

import "errors"

var (
	ErrNotFoundCommand = errors.New("кажется, бот не поддерживает такую команду")
)
