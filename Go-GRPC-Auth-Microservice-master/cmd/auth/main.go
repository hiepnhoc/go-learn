package main

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/entc/integration/migrate/entv2/migrate"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"

	atlas "ariga.io/atlas/sql/migrate"
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/server"
	jaegerTracer "github.com/AleksK1NG/auth-microservice/pkg/jaeger"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/AleksK1NG/auth-microservice/pkg/postgres"
	"github.com/AleksK1NG/auth-microservice/pkg/redis"
	"github.com/AleksK1NG/auth-microservice/pkg/utils"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Starting auth microservice")

	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer psqlDB.Close()

	//migrate db
	//m, err := migrate.New(
	//	"file://migrations",
	//	"postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := m.Up(); err != nil {
	//	log.Fatal(err)
	//}
	dir, err := atlas.NewLocalDir("./migrations/db")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	//if len(os.Args) != 2 {
	//	log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	//}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(context.Background(), "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable", "create_users", opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	tracer, closer, err := jaegerTracer.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	authServer := server.NewAuthServer(appLogger, cfg, psqlDB, redisClient)
	appLogger.Fatal(authServer.Run())
}
