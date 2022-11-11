package main

import (
	"github.com/codfrm/cago/internal/cmd"
	"github.com/codfrm/cago/internal/cmd/gen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cago",
	}

	init := cmd.NewInitCmd()
	rootCmd.AddCommand(init.Commands()...)

	genCmd := gen.NewGenCmd()
	rootCmd.AddCommand(genCmd.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
