package caseservice

import (
	"fmt"

	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/models"
	"github.com/pkg/errors"
)

func DecoError(err error, msg string) error {
	if err != nil {
		return errors.Wrap(err, msg)
	}
	return err
}

type CaseServiceLocalImpl struct {
	Store casedb.CaseStorageService
}

func NewCaseService(store casedb.CaseStorageService) CaseService {
	return CaseServiceLocalImpl{Store: store}
}

func (c CaseServiceLocalImpl) NewCase() (models.Identifier, error) {
	id := models.NewCaseRecordUUID()
	err := c.Store.SaveNewCase(id)
	err = DecoError(err, "Error while storing a new case")
	return id, err
}

func (c CaseServiceLocalImpl) GetCases(spec models.CaseRecordSpec) ([]models.CaseRecord, error) {
	rec, err := c.Store.FindAllCases(spec)
	err = DecoError(err, "Error while fetching all cases")
	return rec, err
}

func (c CaseServiceLocalImpl) FindCase(id string, spec models.CaseRecordSpec) (models.CaseRecord, error) {
	rec, err := c.Store.FindCase(id, spec)
	err = DecoError(err, fmt.Sprintf("Error while fetching case %s", id))
	return rec, err
}

func (c CaseServiceLocalImpl) GetActionRecords(caseId models.Identifier, spec models.CaseActionSpec) ([]models.CaseAction, error) {
	records, err := c.Store.GetContextForCase(caseId)
	err = DecoError(err, "Failed to fetch the case records")
	return records, err
}

func (c CaseServiceLocalImpl) UpdateCase(caseRecord models.CaseRecord) error {
	err := c.Store.UpdateCase(&caseRecord)
	return err
}

func (c CaseServiceLocalImpl) SaveCaseAction(caseAction models.CaseAction) (models.Identifier, error) {
	uuid := models.NewCaseActionUUID()
	err := c.Store.SaveCaseContext(&caseAction)
	return uuid, err
}

func (c CaseServiceLocalImpl) IsActionSupported(action string) bool {
	_, ok := models.BaseSupportedActions[action]
	return ok
}
