package main

//go:generate go run generate.go

import (
	//"sb.im/gosd/cmd"
	"sb.im/gosd/app/cmd"
)

func main() {
	cmd.Execute()
}
