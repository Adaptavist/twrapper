package cmd

import (
	"github.com/google/uuid"
	"log"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func fatalIfNotNil(err error, errStr string) {
	if err != nil {
		log.Fatalf(errStr, err.Error())
	}
}
