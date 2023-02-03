package grpc

import (
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/account/commands"
	"acbs.com.vn/account-service/internal/account/queries"
	"acbs.com.vn/account-service/internal/account/service"
	"acbs.com.vn/account-service/internal/models"
	"acbs.com.vn/account-service/pkg/logger"
	accountService "acbs.com.vn/account-service/proto"
	"context"
	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	ps  *service.AccountService
	//metrics *metrics.ReaderServiceMetrics
}

func NewGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *service.AccountService) *grpcService {
	return &grpcService{log: log, cfg: cfg, v: v, ps: ps}
}

func (s *grpcService) CreateAccount(ctx context.Context, req *accountService.CreateAccountReq) (*accountService.CreateAccountRes, error) {
	//s.metrics.CreateProductGrpcRequests.Inc()
	//
	//ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.CreateProduct")
	//defer span.Finish()

	command := commands.NewCreateAccountCommand(req.GetName())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	account, err := s.ps.Commands.CreateAccount.Handle(ctx, command)
	if err != nil {
		s.log.WarnMsg("CreateProduct.Handle", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	//s.metrics.SuccessGrpcRequests.Inc()
	return &accountService.CreateAccountRes{Name: account.Name}, nil
}

func (s *grpcService) GetAccountByName(ctx context.Context, req *accountService.GetAccountByNameReq) (*accountService.GetAccountByNameRes, error) {
	//s.metrics.GetProductByIdGrpcRequests.Inc()
	//
	//ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetProductById")
	//defer span.Finish()

	query := queries.NewGetAccountByNameQuery(req.GetName())
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	accounts, err := s.ps.Queries.GetAccountByName.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetProductById.Handle", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	//s.metrics.SuccessGrpcRequests.Inc()
	list := models.AccountListToAccounts(accounts)

	return &accountService.GetAccountByNameRes{Accounts: list}, nil
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	//s.metrics.ErrorGrpcRequests.Inc()
	return status.Error(c, err.Error())
}
