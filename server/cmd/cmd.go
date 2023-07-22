package cmd

import (
	"context"

	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/gogo"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Cmd(callback func(ctx context.Context, cmd *cobra.Command) error) cago.FuncComponentCancel {
	return func(ctx context.Context, cancel context.CancelFunc, cfg *configs.Config) error {
		defer cancel()
		rootCmd := &cobra.Command{
			Use: cfg.AppName,
		}
		if err := callback(ctx, rootCmd); err != nil {
			return err
		}
		_ = gogo.Go(func(ctx context.Context) error {
			defer cancel()
			if err := rootCmd.ExecuteContext(ctx); err != nil {
				logger.Ctx(ctx).Error("cmd execute err: %v", zap.Error(err))
				return err
			}
			return nil
		})
		return nil
	}
}
