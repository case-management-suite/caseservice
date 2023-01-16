package caseservice

import (
	"fmt"
	"reflect"

	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/factory"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/scheduler"
)

type CaseServiceClientFactory func(config.AppConfig) server.Server[CaseServiceClient]

type CaseServiceFactory func(casedb.CaseStorageService) CaseService

type CaseServiceServerFactory func(config.AppConfig, CaseService, scheduler.WorkScheduler) server.Server[*CaseServiceGRPCServer]

type CaseMicroervice = server.Server[*CaseServiceGRPCServer]
type CaseClient = server.Server[CaseServiceClient]

type CaseServiceFactories struct {
	factory.FactorySet
	WorkSchedulerFactories    scheduler.WorkSchedulerFactories
	CaseStorageServiceFactory casedb.CaseStorageServiceFactory
	CaseServiceFactory        CaseServiceFactory
	CaseServiceServerFactory  CaseServiceServerFactory
}

func (f CaseServiceFactories) BuildCaseService(appConfig config.AppConfig) (*CaseMicroervice, error) {
	if err := factory.ValidateFactorySet(f); err != nil {
		return nil, fmt.Errorf("factory: %s -> %w;", reflect.TypeOf(f).Name(), err)
	}
	workScheduler, err := f.WorkSchedulerFactories.BuildWorkScheduler(appConfig)
	if err != nil {
		return nil, err
	}

	r := f.CaseServiceServerFactory(
		appConfig,
		f.CaseServiceFactory(
			f.CaseStorageServiceFactory(appConfig.CasesStorage),
		),
		*workScheduler,
	)
	return &r, nil
}
