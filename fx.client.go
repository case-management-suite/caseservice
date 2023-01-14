package caseservice

import (
	"context"

	"github.com/case-management-suite/common/config"
	"go.uber.org/fx"
)

type CaseServiceClientParams struct {
	fx.In
	AppConfig config.AppConfig
}

func newFXCaseServiceClient(lc fx.Lifecycle, params CaseServiceClientParams) CaseService {
	client := NewCaseServiceClient(params.AppConfig)
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return client.Connect()
		},
		OnStop: func(ctx context.Context) error {
			return client.Close()
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
