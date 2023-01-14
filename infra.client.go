package caseservice

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/case-management-suite/api/caseservice/pb"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/models"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func newGRPCConnection(csConf config.CasesServiceConfig) (*grpc.ClientConn, error) {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", csConf.Host, csConf.Port), opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to gRPC client")
	} else {
		log.Info().Str("host", csConf.Host).Int16("port", csConf.Port).Msg("Started gRPC client")
	}
	return conn, err
}

// Implements cases.CaseService
type CaseServiceClient struct {
	Config     config.CasesServiceConfig
	InstanceID int
	clientData *ClientData
}

type ClientData struct {
	IsConnected bool
	connection  *grpc.ClientConn
	client      pb.CaseServiceAPIClient
}

func NewCaseServiceClient(appConfig config.AppConfig) CaseServiceClient {
	cd := ClientData{}
	return CaseServiceClient{Config: appConfig.CasesService, InstanceID: rand.Intn(100), clientData: &cd}
}

func (csc CaseServiceClient) Connect() error {
	conn, err := newGRPCConnection(csc.Config)
	if err != nil {
		log.Error().Err(err).Msg("Client failed to connect to the CasesService")
		return err
	}

	client := pb.NewCaseServiceAPIClient(conn)
	csc.clientData.client = client
	csc.clientData.IsConnected = true
	csc.clientData.connection = conn
	return nil
}

func (csc *CaseServiceClient) Close() error {
	if csc.clientData.IsConnected {
		csc.clientData.connection.Close()
	}
	return nil
}

func (csc CaseServiceClient) checkConnection() error {
	if !csc.clientData.IsConnected {
		return errors.New("CaseServiceClient is not connected. Call connect() before making requests")
	}
	return nil
}

func (c CaseServiceClient) NewCase() (models.Identifier, error) {
	err := c.checkConnection()
	if err != nil {
		return "", err
	}
	client := c.clientData.client
	uuid, err := client.NewCase(context.TODO(), &pb.NewCaseRequest{})
	DecoError(err, "failed to create case")
	if err != nil {
		return "", err
	}
	log.Debug().Str("UUID", uuid.UUID).Msg("Client request for new case is succesful")

	return uuid.UUID, nil
}

func (c CaseServiceClient) GetCases(spec models.CaseRecordSpec) ([]models.CaseRecord, error) {
	err := c.checkConnection()
	if err != nil {
		return []models.CaseRecord{}, err
	}
	client := c.clientData.client
	pbspec := ParseSpec(spec)
	req := pb.FindCasesRequest{Spec: &pbspec}
	resp, err := client.FindCases(context.TODO(), &req)
	err = DecoError(err, "Error while fetching all cases")
	if err != nil {
		return []models.CaseRecord{}, err
	}

	rec := []models.CaseRecord{}

	for _, cr := range resp.CaseRecords {
		rec = append(rec, ParseCaseRecordPB(cr))
	}

	return rec, err
}

func (c CaseServiceClient) FindCase(id string, spec models.CaseRecordSpec) (models.CaseRecord, error) {
	err := c.checkConnection()
	if err != nil {
		return models.CaseRecord{}, err
	}
	client := c.clientData.client
	pbspec := ParseSpec(spec)
	resp, err := client.FindCase(context.TODO(), &pb.FindCaseRequest{UUID: id, Spec: &pbspec})
	err = DecoError(err, "Error while fetching the case")
	if err != nil {
		return models.CaseRecord{}, err
	}
	return ParseCaseRecordPB(resp.CaseRecord), err
}

func (c CaseServiceClient) UpdateCase(record models.CaseRecord) error {
	err := c.checkConnection()
	if err != nil {
		return err
	}
	client := c.clientData.client
	_, err = client.UpdateCase(context.TODO(), &pb.UpdateCaseRequest{
		CaseRecord: ParseCaseRecordModel(record),
	})
	return err
}

func (c CaseServiceClient) GetActionRecords(caseId models.Identifier, spec models.CaseActionSpec) ([]models.CaseAction, error) {
	// records, err := c.Store.GetContextForCase(caseId)
	// err = cases.DecoError(err, "Failed to fetch the case records")
	// return records, err
	return nil, status.Errorf(codes.Unimplemented, "method FindCaseActions not implemented")
}

func (c CaseServiceClient) IsActionSupported(action string) bool {
	_, ok := models.BaseSupportedActions[action]
	return ok
}

func (c CaseServiceClient) SaveCaseAction(action models.CaseAction) (models.Identifier, error) {
	return "", nil
}
