package caseservice

import (
	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/scheduler"
	"go.uber.org/fx"
)

func FxOpts(appConfig config.AppConfig) fx.Option {
	return fx.Options(
		config.FxConfig(appConfig),
		fx.Module("case_service_grpc",
			casedb.NewFxCaseDBService(),
			NewFxCaseService(),
			// CaseServiceModule,
			// rulesengineservice.RulesServiceClientModule,

			fx.Provide(
				func(appConfig config.AppConfig) config.QueueConnectionConfig {
					return appConfig.RulesServiceConfig.QueueConfig
				},
				queue.QueueServiceFactory(appConfig.RulesServiceConfig.QueueType),
				scheduler.NewWorkScheduler,
				NewAPIServerFX,
			),
			fx.Invoke(func(*CaseServiceGRPCServer) {}),
		),
	)

}

func NewCaseServiceGRPCServer(appConfig config.AppConfig) *fx.App {
	return fx.New(
		FxOpts(appConfig),
	)
}
