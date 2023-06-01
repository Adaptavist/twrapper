package utils

import "log"

// Fatal is a wrapper for log.Fatal, if err isn't nil
func FatalIfNotNil(err error, errStr string) {
	if err != nil {
		log.Fatalf(errStr, err.Error())
	}
}

func FatalIfNil(obj interface{}, errStr string) {
	if obj == nil {
		log.Fatalf(errStr)
	}
}

func FatalIfNotOk(ok bool, errStr string) {
	if !ok {
		log.Fatalf(errStr)
	}
}

// DieOnError will print error and exit if it is not nil
func DieOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FatalIfEmpty(str, errStr string) {
	if len(str) == 0 {
		log.Fatal(errStr)
	}
}

func FatalIfFalse(b bool, errStr string) {
	if !b {
		log.Fatal(errStr)
	}
}
