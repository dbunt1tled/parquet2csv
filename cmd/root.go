package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals // need for init commands
	Use:   "converter",
	Short: "Converter CLI",
	Long:  "Converter parquet to csv",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) //nolint:forbidigo // print error
		os.Exit(1)
	}
}
