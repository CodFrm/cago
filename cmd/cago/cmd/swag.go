package cmd

import (
	"github.com/codfrm/cago/third_party/swag/gen"
	"github.com/spf13/cobra"
)

type swagCmd struct {
	dir string
}

func NewSwagCmd() *swagCmd {
	return &swagCmd{}
}

func (s *swagCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "swag",
		Short: "生成swagger文档",
		RunE:  s.gen,
	}
	ret.Flags().StringVarP(&s.dir, "dir", "d", "./", "搜索目录")
	return []*cobra.Command{ret}
}

func (s *swagCmd) gen(cmd *cobra.Command, args []string) error {
	return gen.New().Build(&gen.Config{
		SearchDir:           s.dir,
		Excludes:            "",
		MainAPIFile:         "main.go",
		PropNamingStrategy:  "camelcase",
		OutputDir:           "./docs",
		ParseVendor:         false,
		ParseDependency:     false,
		MarkdownFilesDir:    "",
		ParseInternal:       false,
		GeneratedTime:       false,
		CodeExampleFilesDir: "",
		ParseDepth:          100,
	})
}
