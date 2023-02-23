package server

import (
	sessRepository "github.com/AleksK1NG/auth-microservice/internal/session/repository"
	sessUseCase "github.com/AleksK1NG/auth-microservice/internal/session/usecase"
	"github.com/AleksK1NG/auth-microservice/internal/user/delivery/darp"
	authServerGRPC "github.com/AleksK1NG/auth-microservice/internal/user/delivery/grpc/service"
	userRepository "github.com/AleksK1NG/auth-microservice/internal/user/repository"
	userUseCase "github.com/AleksK1NG/auth-microservice/internal/user/usecase"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	pb "github.com/dapr/go-sdk/dapr/proto/runtime/v1"
	"github.com/dapr/go-sdk/service/common"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/interceptors"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/AleksK1NG/auth-microservice/pkg/metric"
)

var (
	service common.Service
)

// GRPC Auth Server
type Server struct {
	logger      logger.Logger
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
}

// Server constructor
func NewAuthServer(logger logger.Logger, cfg *config.Config, db *sqlx.DB, redisClient *redis.Client) *Server {
	return &Server{logger: logger, cfg: cfg, db: db, redisClient: redisClient}
}

// Run service
func (s *Server) Run() error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	im := interceptors.NewInterceptorManager(s.logger, s.cfg, metrics)
	userRepo := userRepository.NewUserPGRepository(s.db)
	sessRepo := sessRepository.NewSessionRepository(s.redisClient, s.cfg)
	userRedisRepo := userRepository.NewUserRedisRepo(s.redisClient, s.logger)
	userUC := userUseCase.NewUserUseCase(s.logger, userRepo, userRedisRepo)
	sessUC := sessUseCase.NewSessionUseCase(sessRepo, s.cfg)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	serverGrpc := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Server.Timeout * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(serverGrpc)
	}

	authGRPCServer := authServerGRPC.NewAuthServerGRPC(s.logger, s.cfg, userUC, sessUC)
	userService.RegisterUserServiceServer(serverGrpc, authGRPCServer)

	pb.RegisterAppCallbackServer(serverGrpc, darp.NewDarpGprc(s.logger, s.cfg, authGRPCServer))

	grpc_prometheus.Register(serverGrpc)
	//http.Handle("/metrics", promhttp.Handler())

	//if service, err = daprd.NewService(":9000"); err != nil {
	//	//log.Infof("failed to start the server: %v", err)
	//	return err
	//}

	//darphttp := darp.NewDarpHttp(s.logger, s.cfg, authGRPCServer)
	//
	//darphttp.AddInvocationHandler(s.logger, service, "/register")

	go func() {
		s.logger.Infof("Server is listening on port: %v", s.cfg.Server.Port)
		if err := serverGrpc.Serve(l); err != nil {
			s.logger.Fatal(err)
		}
	}()

	//if err := service.Start(); err != nil && err != http.ErrServerClosed {
	//	s.logger.Fatal("error: %v", err)
	//}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	serverGrpc.GracefulStop()
	s.logger.Info("Server Exited Properly")

	return nil
}
