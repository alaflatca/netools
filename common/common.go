package common

import "os"

func FileExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}
