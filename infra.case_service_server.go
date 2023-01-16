package caseservice

import (
	"context"
	"fmt"
	"net"

	"github.com/case-management-suite/api/caseservice/pb"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/ctxutils"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/models"
	"github.com/case-management-suite/scheduler"
	"github.com/rs/zerolog"
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
	server.ServerUtils
}

func (s CaseServiceServer) NewCase(context.Context, *pb.NewCaseRequest) (*pb.UUIDResponse, error) {
	s.Logger.Debug().Str("service", "grpc").Msg("NewCase")
	id, err := s.casesService.NewCase()

	if err != nil {
		return &pb.UUIDResponse{}, err
	}
	s.Logger.Debug().Str("UUID", id).Msg("Saved ID")
	return &pb.UUIDResponse{UUID: id}, err
}
func (s CaseServiceServer) FindCase(ctx context.Context, req *pb.FindCaseRequest) (*pb.FindCaseResponse, error) {
	s.Logger.Debug().Str("service", "grpc").Str("UUID", req.UUID).Msg("FindCase")
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
	s.Logger.Debug().Str("service", "grpc").Str("UUID", req.CaseRecord.ID).Msg("UpdateCase")
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

func (c *CaseServiceGRPCServer) GetName() string {
	return "Case Service gRPC Server"
}

func (c *CaseServiceGRPCServer) GetServerConfig() *server.ServerConfig {
	return &server.ServerConfig{
		Type: server.GRPCServerType,
		Port: int(c.Port),
	}
}

func (c *CaseServiceGRPCServer) Start(ctx context.Context) error {
	ctx = ctxutils.DecorateContext(ctx, ctxutils.ContextDecoration{Name: "CaseServiceGRPCService"})
	logger := zerolog.Ctx(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to listen")
	}

	c.GrpcServer = grpc.NewServer()

	// attach the Ping service to the server
	pb.RegisterCaseServiceAPIServer(c.GrpcServer, c.CaseServer)

	if err := c.Scheduler.Start(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to start scheduler")
		return err
	}

	go func() {
		c.CaseServer.Logger.Info().Int16("port", c.Port).Msg("gRPC service for cases started")
		if err := c.GrpcServer.Serve(lis); err != nil {
			c.CaseServer.Logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	newCtx := context.Background()
	c.CaseServer.ListenForCaseUpdates(newCtx)

	return nil
}

func (c *CaseServiceGRPCServer) Stop(_ context.Context) error {
	c.GrpcServer.Stop()
	return c.Scheduler.Stop()
}

func NewCaseServiceAPIServer(appConfig config.AppConfig, caseService CaseService, wshservice scheduler.WorkScheduler) server.Server[*CaseServiceGRPCServer] {
	port := appConfig.CasesService.Port
	return server.NewServer(func(su server.ServerUtils) *CaseServiceGRPCServer {
		caseServer := CaseServiceServer{casesService: caseService, Scheduler: wshservice, ServerUtils: su}
		grpcServer := CaseServiceGRPCServer{Port: port, CaseService: caseService, Scheduler: wshservice, CaseServer: caseServer}
		return &grpcServer
	}, appConfig)
}

func _() CaseServiceServerFactory {
	return NewCaseServiceAPIServer
}
