package caseservice

import (
	"context"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/server"
	"go.uber.org/fx"
)

type CaseServiceClientParams struct {
	fx.In
	AppConfig config.AppConfig
}

func newFXCaseServiceClient(lc fx.Lifecycle, params CaseServiceClientParams) server.Server[CaseServiceClient] {
	client := NewCaseServiceClient(params.AppConfig)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return client.Stop(ctx)
		},
	})
	return client
}

func NewFXCaseServiceClient() fx.Option {
	return fx.Module("case_service_client",
		fx.Provide(
			newFXCaseServiceClient,
		),
	)
}
