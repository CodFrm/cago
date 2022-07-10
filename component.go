package cago

import (
	"context"

	"github.com/codfrm/cago/config"
)

type Component interface {
	Start(ctx context.Context, cfg *config.Config) error
	CloseHandle()
}

type ComponentCancel interface {
	Component
	StartCancel(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) error
	CloseHandle()
}

type FuncComponent func(ctx context.Context, cfg *config.Config) error

func (f FuncComponent) Start(ctx context.Context, cfg *config.Config) error {
	return f(ctx, cfg)
}

func (f FuncComponent) CloseHandle() {
}

type FuncComponentCancel func(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) error

func (f FuncComponentCancel) Start(ctx context.Context, cfg *config.Config) error {
	return f.StartCancel(ctx, nil, cfg)
}

func (f FuncComponentCancel) StartCancel(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) error {
	return f(ctx, cancel, cfg)
}

func (f FuncComponentCancel) CloseHandle() {
}
