package cmd

import (
	"sb.im/gosd/config"
	"sb.im/gosd/database"
	"sb.im/gosd/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mustGetString(flags *pflag.FlagSet, flag string) string {
	s, err := flags.GetString(flag)
	checkErr(err)
	return s
}

func mustGetBool(flags *pflag.FlagSet, flag string) bool {
	b, err := flags.GetBool(flag)
	checkErr(err)
	return b
}

func mustGetUint(flags *pflag.FlagSet, flag string) uint {
	b, err := flags.GetUint(flag)
	checkErr(err)
	return b
}

type cobraFunc func(cmd *cobra.Command, args []string)
type exCobraFunc func(cmd *cobra.Command, args []string, data exData)

type exConfig struct{}

type exData struct {
	store *storage.Storage
}

func ex(fn exCobraFunc, cfg exConfig) cobraFunc {
	return func(cmd *cobra.Command, args []string) {
		parse := config.NewParser()
		opts, err := parse.ParseEnvironmentVariables()
		if err != nil {
			panic(err)
		}

		db, err := database.NewConnectionPool(
			opts.DatabaseURL(),
			opts.DatabaseMinConns(),
			opts.DatabaseMaxConns(),
		)

		data := exData{
			store: storage.NewStorage(db),
		}

		fn(cmd, args, data)
	}
}
