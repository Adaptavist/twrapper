package utils

import (
	"os"
)

// handleChdirFlag will chdir to the given path when the flag is set
func ChangeDirOrDie(dir string) {
	err := os.Chdir(dir)
	FatalIfNotNil(err, "failed to chdir to %s")
}
