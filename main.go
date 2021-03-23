package main

//go:generate go run generate.go

import (
	"sb.im/gosd/cmd"
)

func main() {
	cmd.Execute()
}
