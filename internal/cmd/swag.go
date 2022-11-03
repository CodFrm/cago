package cmd

import (
	"github.com/codfrm/cago/third_party/swag/gen"
	"github.com/spf13/cobra"
	"github.com/swaggo/swag/format"
)

type SwagCmd struct {
	dir string
}

func NewSwagCmd() *SwagCmd {
	return &SwagCmd{}
}

func (s *SwagCmd) Commands() []*cobra.Command {
	ret := &cobra.Command{
		Use:   "swag",
		Short: "生成swagger文档",
		RunE:  s.gen,
	}
	fmt := &cobra.Command{
		Use:   "fmt",
		Short: "格式化swagger注释",
		RunE:  s.fmt,
	}
	ret.AddCommand(fmt)
	ret.Flags().StringVarP(&s.dir, "dir", "d", "./", "搜索目录")
	fmt.Flags().StringVarP(&s.dir, "dir", "d", "./", "搜索目录")
	return []*cobra.Command{ret}
}

func (s *SwagCmd) gen(cmd *cobra.Command, args []string) error {
	if err := s.fmt(cmd, args); err != nil {
		return err
	}
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

func (s *SwagCmd) fmt(cmd *cobra.Command, args []string) error {
	return format.New().Build(&format.Config{
		SearchDir: s.dir,
		Excludes:  "",
		MainFile:  "main.go",
	})
}
