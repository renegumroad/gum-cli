package utils

import "github.com/renegumroad/gum-cli/internal/log"

func CheckFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
