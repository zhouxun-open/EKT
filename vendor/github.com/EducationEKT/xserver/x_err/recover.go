package x_err

import "fmt"

func Recover() {
	if err := recover(); err != nil {
		fmt.Errorf("recovered an error")
	}
}
