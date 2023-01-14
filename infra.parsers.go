package caseservice

import (
	"github.com/case-management-suite/api/caseservice/pb"
	"github.com/case-management-suite/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func ParseCaseRecordModel(cr models.CaseRecord) *pb.CaseRecord {
	var deletedTime *timestamppb.Timestamp = nil
	if cr.DeletedAt.Valid {
		deletedTime = timestamppb.New(cr.DeletedAt.Time)
	}
	return &pb.CaseRecord{
		ID:        cr.ID,
		CreatedAt: timestamppb.New(cr.CreatedAt),
		UpdatedAt: timestamppb.New(cr.UpdatedAt),
		DeletedAt: deletedTime,
		Status:    cr.Status,
	}
}

func ParseCaseRecordPB(pbrec *pb.CaseRecord) models.CaseRecord {
	var deletedAt gorm.DeletedAt

	if pbrec.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: pbrec.DeletedAt.AsTime(), Valid: true}
	} else {
		deletedAt = gorm.DeletedAt{Time: pbrec.DeletedAt.AsTime(), Valid: false}
	}

	return models.CaseRecord{
		ID:        pbrec.ID,
		Status:    pbrec.Status,
		CreatedAt: pbrec.CreatedAt.AsTime(),
		UpdatedAt: pbrec.UpdatedAt.AsTime(),
		DeletedAt: deletedAt,
	}
}

func ParseSpec(spec models.CaseRecordSpec) pb.Spec {
	return pb.Spec{Fields: spec.GetMap()}
}
