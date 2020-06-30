package main

//go:generate go run generate.go

import (
	"sb.im/gosd/cli"
)

func main() {
	cli.Parse()
}
