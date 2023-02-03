package server

import (
	"acbs.com.vn/account-service/pkg/constants"
	"context"
	"github.com/heptiolabs/healthcheck"
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
