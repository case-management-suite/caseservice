package caseservice

import (
	// "dig"
	"context"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/scheduler"
	"go.uber.org/fx"
)

type CaseServiceGRPCServerModuleParams struct {
	fx.In
	CaseService CaseService
	Scheduler   scheduler.WorkScheduler
	AppConfig   config.AppConfig
}

func NewAPIServerFX(lc fx.Lifecycle, params CaseServiceGRPCServerModuleParams) *CaseServiceGRPCServer {

	server := NewCaseServiceAPIServer(params.AppConfig, params.CaseService, params.Scheduler)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// serve API
			err := server.Start(ctx)
			// start the server
			return err
		},
		OnStop: func(_ context.Context) error {
			// return http.
			return server.Stop()
		},
	})

	return &server
}
