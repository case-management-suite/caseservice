package caseservice_test

import (
	"log"
	"testing"

	"github.com/case-management-suite/casedb/mocks"
	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/models"
	"github.com/golang/mock/gomock"
)

func TestNewCase(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()
	db := mocks.NewMockCaseStorageService(ctrl)
	db.EXPECT().SaveNewCase(gomock.Any()).DoAndReturn(func(id models.Identifier) error {
		log.Printf("Captured %v", id)
		return nil
	})

	ctrler := caseservice.NewCaseService(db)
	id, err := ctrler.NewCase()

	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}

	log.Printf("ID: %v", id)
}

func TestGetCases(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()
	db := mocks.NewMockCaseStorageService(ctrl)
	db.EXPECT().FindAllCases(gomock.Any()).Return([]models.CaseRecord{{ID: "id1", Status: "NEW_CASE"}, {ID: "id2", Status: "NEW_CASE"}}, nil)

	ctrler := caseservice.NewCaseService(db)
	id, err := ctrler.GetCases(models.NewCaseRecordSpec(true))

	if err != nil {
		t.Fatalf("There was an error: %v", err)
	}

	if len(id) != 2 {
		t.Fatalf("Expected two models in the response: %v", err)
	}

	log.Printf("ID: %v", id)
}

// func TestExecuteAction(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	// Assert that Bar() is invoked.
// 	defer ctrl.Finish()
// 	db := mocks.NewMockCaseStorageService(ctrl)
// 	db.EXPECT().FindCase(gomock.Any(), gomock.Any()).Return(models.CaseRecord{ID: "id1", Status: cases.CaseStatusList.NewCase}, nil)

// 	db.EXPECT().UpdateCase(gomock.Any()).DoAndReturn(func(model *models.CaseRecord) error {
// 		log.Printf("Captured %v", model)
// 		if model.Status != cases.CaseStatusList.Started {
// 			t.Fatalf("Expected an update status to STARTED, but was %v", model.Status)
// 		}
// 		return nil
// 	})

// 	db.EXPECT().SaveCaseContext(gomock.Any()).DoAndReturn(func(model *models.CaseAction) error {
// 		log.Printf("Captured %v", model)
// 		if model.CaseRecord.Status != cases.CaseStatusList.Started {
// 			t.Fatalf("Expected an update status to STARTED, but was %v", model.CaseRecord.Status)
// 		}
// 		return nil
// 	})

// 	ctrler := cases.NewCaseService(db, cases.NewRulesEngineService())
// 	record, err := ctrler.ExecuteAction("id1", "START")

// 	if err != nil {
// 		t.Fatalf("There was an error: %v", err)
// 	}

// 	if record.Status != cases.CaseStatusList.Started {
// 		t.Fatalf("The code was expected to be Started but was %v", record.Status)
// 	}

// 	log.Printf("record: %v", record)
// }
