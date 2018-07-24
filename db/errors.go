package db

import "errors"

var (
	NoSuchKeyError   = errors.New("no such key in database")
	InvalidTypeError = errors.New("invalid result type")
)
