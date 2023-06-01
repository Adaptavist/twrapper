package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Dexit dumps the output as yaml and exists
func Dexit(obj interface{}) {
	out, _ := yaml.Marshal(obj)
	fmt.Println(string(out))
	os.Exit(1)
}
