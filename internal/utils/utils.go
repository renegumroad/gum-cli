package utils

import "log"

func CheckFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
