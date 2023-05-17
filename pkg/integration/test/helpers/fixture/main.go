package fixture

import (
	"embed"
	"fmt"
	"os"
)

//go:embed fixtures/*
//go:embed fixtures/random
var Fixtures embed.FS

func Get(f string) []byte {
	b, err := Fixtures.ReadFile(fmt.Sprintf("fixtures/%s", f))
	if err != nil {
		panic(err)
	}
	return b
}

func Write(fixture string, destination string) {
	if err := os.WriteFile(destination, Get(fixture), 0644); err != nil {
		panic(err)
	}
}
