package stat

import "errors"

var (
	WaitingSyncErr = errors.New("waiting synchronize blockchain")
	NoAccountErr   = errors.New("no such address")
)
