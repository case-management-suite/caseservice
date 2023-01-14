package caseservice

import (
	"context"
	"fmt"
	"net"

	"github.com/case-management-suite/api/caseservice/pb"
	"github.com/case-management-suite/common"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/ctxutils"
	"github.com/case-management-suite/models"
	"github.com/case-management-suite/scheduler"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CaseServiceServer represents the gRPC server
type CaseServiceServer struct {
	*pb.UnimplementedCaseServiceAPIServer
	casesService CaseService
	// Rules        rulesengineservice.RulesServiceClient
	Scheduler scheduler.WorkScheduler
	common.ServerUtils
}

func (s CaseServiceServer) NewCase(context.Context, *pb.NewCaseRequest) (*pb.UUIDResponse, error) {
	log.Debug().Str("service", "grpc").Msg("NewCase")
	id, err := s.casesService.NewCase()

	if err != nil {
		return &pb.UUIDResponse{}, err
	}
	log.Debug().Str("UUID", id).Msg("Saved ID")
	return &pb.UUIDResponse{UUID: id}, err
}
func (s CaseServiceServer) FindCase(ctx context.Context, req *pb.FindCaseRequest) (*pb.FindCaseResponse, error) {
	log.Debug().Str("service", "grpc").Str("UUID", req.UUID).Msg("FindCase")
	spec := models.NewCaseRecordSpec(false)
	for k, v := range req.Spec.Fields {
		err := spec.Set(k, v)
		if err != nil {
			return &pb.FindCaseResponse{}, err
		}
	}
	rec, err := s.casesService.FindCase(req.UUID, spec)
	if err != nil {
		return &pb.FindCaseResponse{}, err
	}
	return &pb.FindCaseResponse{CaseRecord: ParseCaseRecordModel(rec)}, nil
}
func (s CaseServiceServer) FindCases(_ context.Context, req *pb.FindCasesRequest) (*pb.FindCasesResponse, error) {
	spec := models.NewCaseRecordSpec(false)
	for k, v := range req.Spec.Fields {
		err := spec.Set(k, v)
		if err != nil {
			return &pb.FindCasesResponse{}, err
		}
	}
	caseList, err := s.casesService.GetCases(spec)
	if err != nil {
		return &pb.FindCasesResponse{}, err
	}
	results := []*pb.CaseRecord{}
	for _, cr := range caseList {
		results = append(results, ParseCaseRecordModel(cr))
	}
	return &pb.FindCasesResponse{CaseRecords: results}, nil
}
func (CaseServiceServer) FindCaseActions(context.Context, *pb.FindCaseActionsRequest) (*pb.FindCaseActionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindCaseActions not implemented")
}

func (s CaseServiceServer) UpdateCase(ctx context.Context, req *pb.UpdateCaseRequest) (*pb.UpdateCaseResponse, error) {
	// id, err := s.casesService.UpdateCase(ParseCaseRecordPB(req.CaseRecord))
	log.Debug().Str("service", "grpc").Str("UUID", req.CaseRecord.ID).Msg("UpdateCase")
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCaseRequest not implemented")
}

func (s CaseServiceServer) ListenForCaseUpdates(ctx context.Context) error {
	err := s.Scheduler.ListenForCaseUpdates(func(caseRecord models.CaseRecord) error {
		return s.casesService.UpdateCase(caseRecord)
	}, ctx)
	return err
}

func _() pb.CaseServiceAPIServer {
	s := CaseServiceServer{}
	return s
}

type CaseServiceGRPCServer struct {
	CaseService CaseService
	Scheduler   scheduler.WorkScheduler
	Port        int16
	CaseServer  CaseServiceServer
	GrpcServer  *grpc.Server
}

func (c *CaseServiceGRPCServer) Start(ctx context.Context) error {
	ctx = ctxutils.DecorateContext(ctx, ctxutils.ContextDecoration{Name: "CaseServiceGRPCService"})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	c.GrpcServer = grpc.NewServer()

	// attach the Ping service to the server
	pb.RegisterCaseServiceAPIServer(c.GrpcServer, c.CaseServer)

	if err := c.Scheduler.Start(ctx); err != nil {
		log.Error().Err(err).Msg("failed to start scheduler")
		return err
	}

	go func() {
		log.Info().Int16("port", c.Port).Msg("gRPC service for cases started")
		if err := c.GrpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	c.CaseServer.ListenForCaseUpdates(ctx)

	return nil
}

func (c *CaseServiceGRPCServer) Stop() error {
	c.GrpcServer.Stop()
	return c.Scheduler.Stop()
}

func NewCaseServiceAPIServer(appConfig config.AppConfig, caseService CaseService, wshservice scheduler.WorkScheduler) CaseServiceGRPCServer {
	port := appConfig.CasesService.Port
	logConf := appConfig.LogConfig
	caseServer := CaseServiceServer{casesService: caseService, Scheduler: wshservice, ServerUtils: common.ServerUtils{Logger: logConf.Logger.Level(logConf.CasesService)}}
	return CaseServiceGRPCServer{Port: port, CaseService: caseService, Scheduler: wshservice, CaseServer: caseServer}
}
