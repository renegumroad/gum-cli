package utils

import "github.com/renehernandez/gum-cli/internal/log"

func CheckFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
