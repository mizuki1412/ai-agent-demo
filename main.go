package main

import (
	"github.com/mizuki1412/go-core-kit/v2/cli"
	"github.com/spf13/cobra"
)

func main() {
	cli.RootCMD(&cobra.Command{
		Use: "main",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
	cli.Execute()
}
