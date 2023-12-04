package utils

import "github.com/lowk3v/dumpsc/config"

func HandleError(err error, message string) bool {
	if err != nil {
		config.Log.Error(message+" ", err)
		return true
	}
	return false
}
