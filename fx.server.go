package caseservice

import (
	// "dig"
	"context"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/scheduler"
	"go.uber.org/fx"
)

type CaseServiceGRPCServerModuleParams struct {
	fx.In
	CaseService CaseService
	Scheduler   scheduler.WorkScheduler
	AppConfig   config.AppConfig
}

func NewAPIServerFX(lc fx.Lifecycle, params CaseServiceGRPCServerModuleParams) *server.Server[*CaseServiceGRPCServer] {
	server := NewCaseServiceAPIServer(params.AppConfig, params.CaseService, params.Scheduler)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := server.Start(ctx)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return server.Stop(ctx)
		},
	})

	return &server
}
