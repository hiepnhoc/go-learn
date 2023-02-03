package server

import (
	"acbs.com.vn/account-service/db"
	"acbs.com.vn/account-service/pkg/constants"
	"context"
	"fmt"
	"github.com/heptiolabs/healthcheck"
	"log"
	"net/http"
	"time"
)

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	health.AddLivenessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck(constants.Postgres, healthcheck.AsyncWithContext(ctx, func() error {
		return s.db.Ping()
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	//health.AddReadinessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
	//	//_, err := s.kafkaConn.Brokers()
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	go func() {
		s.log.Infof("Writer microservice Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) Migrate() {
	log.Println("migrating...")

	dsn := fmt.Sprintf("%s://%s:%d/%s?sslmode=%s&user=%s&password=%s",
		"postgres",
		s.cfg.Postgresql.Host,
		s.cfg.Postgresql.Port,
		s.cfg.Postgresql.DBName,
		s.cfg.Postgresql.SSLMode,
		s.cfg.Postgresql.User,
		s.cfg.Postgresql.Password)

	migrator := db.Migrator(db.WithDSN(dsn))
	// sqlDatabase := db.New(cfg)

	migrator.DB = s.db.DB

	if err := migrator.DB.Ping(); err != nil {
		log.Fatalf("Migrate Ping error %v", err)
	}

	// todo: accept cli flag for other operations
	migrator.Up()

	log.Println("done migration.")
}
