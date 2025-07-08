package utils

import "os"

func IsDev() bool {
	for _, v := range os.Args {
		if v == "--dev" {
			return true
		}
	}
	return false
}
