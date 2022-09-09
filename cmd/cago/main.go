package main

import (
	"github.com/codfrm/cago/cmd/cago/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cago",
	}

	init := cmd.NewInitCmd()
	rootCmd.AddCommand(init.Commands()...)

	swag := cmd.NewSwagCmd()
	rootCmd.AddCommand(swag.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
