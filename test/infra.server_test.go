package caseservice_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/models"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/scheduler"
	"github.com/case-management-suite/testutil"
	"github.com/rs/zerolog"
)

func CaseRecordIDMatches(id string) func(rec models.CaseRecord, err error) bool {
	return func(rec models.CaseRecord, err error) bool { return rec.ID == id }
}

func TestListenForCaseUpdate(t *testing.T) {
	appConfig := config.NewLocalAppConfig()
	// appConfig.CasesStorage.LogSQL = true
	appConfig.RulesServiceConfig.QueueConfig.LogLevel = zerolog.DebugLevel
	storage := casedb.NewSQLCaseStorageService(appConfig.CasesStorage)
	caseService := caseservice.NewCaseService(storage)

	newQueueService := queue.QueueServiceFactory(config.GoChannels)
	qs := newQueueService(appConfig.RulesServiceConfig.QueueConfig, appConfig.LogConfig)
	workservice := scheduler.NewWorkScheduler(qs, appConfig)
	server := caseservice.NewCaseServiceAPIServer(appConfig, caseService, workservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	workservice.Start(ctx)
	defer workservice.Stop()
	// server.Start(ctx)

	err := server.CaseServer.ListenForCaseUpdates(ctx)
	testutil.AssertNilError(err, t)

	out, err := qs.Listen(appConfig.RulesServiceConfig.QueueConfig.CaseNotificationsChannel, ctx)
	testutil.AssertNilError(err, t)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case msg := <-out.Out:
				msg.Ack()
				// TODO: Removethis:
				time.Sleep(time.Second)
				wg.Done()

				return
			case <-time.After(time.Second):
			case <-out.Done:
				return
			}
		}
	}()

	id, err := server.CaseService.NewCase()
	testutil.AssertNilError(err, t)

	record, err := server.CaseService.FindCase(id, models.NewCaseRecordSpec(true))
	testutil.AssertNilError(err, t)
	testutil.AssertEq(id, record.ID, t)

	testutil.AssertEq("NEW_CASE", record.Status, t)

	new_record := models.CaseRecord{ID: record.ID, Status: "STARTED"}

	err = server.Scheduler.NotifyCaseUpdate(new_record, ctx)
	testutil.AssertNilError(err, t)

	wg.Wait()

	record, err = server.CaseService.FindCase(id, models.NewCaseRecordSpec(true))
	testutil.AssertNilError(err, t)
	testutil.AssertEq(id, record.ID, t)

	testutil.AssertEq("STARTED", record.Status, t)
}
