package cmd

import (
	"github.com/spf13/cobra"
)

type initCmd struct {
}

func NewInitCmd() *initCmd {
	return &initCmd{}
}

func (e *initCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "init [name]",
		Short: "初始化项目",
		RunE:  e.exec,
		Args:  cobra.ExactArgs(1),
	}
	return []*cobra.Command{ret}
}

func (e *initCmd) exec(cmd *cobra.Command, args []string) error {

	return nil
}
