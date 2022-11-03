package main

import (
	cmd2 "github.com/codfrm/cago/internal/cmd"
	gen2 "github.com/codfrm/cago/internal/cmd/gen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cago",
	}

	init := cmd2.NewInitCmd()
	rootCmd.AddCommand(init.Commands()...)

	swag := cmd2.NewSwagCmd()
	rootCmd.AddCommand(swag.Commands()...)

	gen := gen2.NewGenCmd()
	rootCmd.AddCommand(gen.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
