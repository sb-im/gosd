package cmd

import (
	"github.com/spf13/pflag"
)

func mustGetUint(flags *pflag.FlagSet, flag string) uint {
	b, _ := flags.GetUint(flag)
	return b
}

func mustGetString(flags *pflag.FlagSet, flag string) string {
	s, _ := flags.GetString(flag)
	return s
}
