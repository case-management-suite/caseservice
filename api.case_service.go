package caseservice

import "github.com/case-management-suite/models"

type CaseService interface {
	NewCase() (models.Identifier, error)
	GetCases(models.CaseRecordSpec) ([]models.CaseRecord, error)
	FindCase(id models.Identifier, spec models.CaseRecordSpec) (models.CaseRecord, error)
	UpdateCase(models.CaseRecord) error
	// ExecuteAction(id models.Identifier, action string) (models.CaseRecord, error)
	SaveCaseAction(models.CaseAction) (models.Identifier, error)
	GetActionRecords(caseId models.Identifier, spec models.CaseActionSpec) ([]models.CaseAction, error)
	IsActionSupported(action string) bool
}
