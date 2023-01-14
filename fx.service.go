package caseservice

import (
	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/common/config"

	"go.uber.org/fx"
)

type CaseServiceParms struct {
	fx.In
	AppConfig          config.AppConfig
	CaseStorageService casedb.CaseStorageService
}

type CaseServiceResult struct {
	fx.Out
	CaseService CaseService
}

func ProvideResult(params CaseServiceParms) CaseServiceResult {
	service := NewCaseService(params.CaseStorageService)
	return CaseServiceResult{CaseService: service}
}

func NewFxCaseService() fx.Option {
	return fx.Module("case_service",
		fx.Provide(
			ProvideResult,
		),
	)
}
