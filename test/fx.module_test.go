package caseservice_test

import (
	"testing"

	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/testutil"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestNewCaseServiceGRPCServer(t *testing.T) {
	testutil.AppFx(t, caseservice.FxOpts(config.NewLocalTestAppConfig()), func(a *fxtest.App) {})
}

func TestNewFXCaseServiceClient(t *testing.T) {
	if err := fx.ValidateApp(caseservice.NewFXCaseServiceClient()); err != nil {
		log.Error().Err(err).Msg("Failed validation")
	}
}

func TestNewFXCaseService(t *testing.T) {
	if err := fx.ValidateApp(caseservice.NewFxCaseService()); err != nil {
		log.Error().Err(err).Msg("Failed validation")
	}
}
