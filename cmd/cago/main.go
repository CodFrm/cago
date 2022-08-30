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

	initCmd := cmd.NewInitCmd()

	rootCmd.AddCommand(initCmd.Commands()...)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
