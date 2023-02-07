package server

import (
	"acbs.com.vn/account-service/client/esign"
	"acbs.com.vn/account-service/config"
	"acbs.com.vn/account-service/internal/account/repository"
	"acbs.com.vn/account-service/internal/account/service"
	"acbs.com.vn/account-service/pkg/interceptors"
	"acbs.com.vn/account-service/pkg/logger"
	"acbs.com.vn/account-service/pkg/postgres"
	"context"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	//kafkaConn *kafka.Conn
	ps *service.AccountService
	im interceptors.InterceptorManager
	db *sqlx.DB

	//metrics   *metrics.WriterServiceMetrics
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.log)
	//s.metrics = metrics.NewWriterServiceMetrics(s.cfg)

	pgxConn, err := postgres.NewSqlx(s.cfg.Postgresql)

	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn")
	}

	s.db = pgxConn
	//s.log.Infof("postgres co-nnected: %v", pgxConn.Stat().TotalConns())
	defer pgxConn.Close()

	//kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	//defer kafkaProducer.Close() // nolint: errcheck

	accountRepo := repository.NewAccountRepository(s.log, s.cfg, s.db)
	s.ps = service.NewAccountService(s.log, s.cfg, accountRepo)
	//productMessageProcessor := kafkaConsumer.NewProductMessageProcessor(s.log, s.cfg, s.v, s.ps, s.metrics)

	//s.log.Info("Starting Writer Kafka consumers")
	//cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	//go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, productMessageProcessor.ProcessMessages)

	closeGrpcServer, grpcServer, err := s.newAccountGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	//if err := s.connectKafkaBrokers(ctx); err != nil {
	//	return errors.Wrap(err, "s.connectKafkaBrokers")
	//}
	//defer s.kafkaConn.Close() // nolint: errcheck
	//
	//if s.cfg.Kafka.InitTopics {
	//	s.initKafkaTopics(ctx)
	//}

	clientRest := resty.New()

	esignClient := esign.NewEsignClient(s.log, s.cfg, clientRest)

	esignClient.GetDetail(context.Background())

	s.Migrate()
	s.runHealthCheck(ctx)
	//s.runMetrics(cancel)

	//if s.cfg.Jaeger.Enable {
	//	tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
	//	if err != nil {
	//		return err
	//	}
	//	defer closer.Close() // nolint: errcheck
	//	opentracing.SetGlobalTracer(tracer)
	//}

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
