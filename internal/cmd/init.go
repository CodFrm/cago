package cmd

import (
	"github.com/spf13/cobra"
)

type InitCmd struct {
}

func NewInitCmd() *InitCmd {
	return &InitCmd{}
}

func (e *InitCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "init [name]",
		Short: "初始化项目",
		RunE:  e.exec,
		Args:  cobra.ExactArgs(1),
	}
	return []*cobra.Command{ret}
}

func (e *InitCmd) exec(cmd *cobra.Command, args []string) error {

	return nil
}
