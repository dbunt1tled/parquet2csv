package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals // need for init commands
	Use:   "csv2parquet",
	Short: "Converter CLI",
	Long:  "Converter parquet ⇄ csv",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) //nolint:forbidigo // print error
		os.Exit(1)
	}
}
