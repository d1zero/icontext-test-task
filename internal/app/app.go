package app

import (
	"context"
	"errors"
	"fmt"
	v1 "icontext-test-task/internal/controller/v1"
	"icontext-test-task/internal/gateway"
	"icontext-test-task/internal/service"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// db driver
	_ "github.com/jackc/pgx/v5/stdlib"

	// gomigrate migration resolver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Run() {
	// loading config
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	// creating logger
	atom := zap.NewAtomicLevel()
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		atom,
	)
	logger := zap.New(zapCore)
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)

	l := logger.Sugar()
	atom.SetLevel(zapcore.Level(cfg.Logger.Level))
	l.Infof("logger initialized successfully")

	// connecting to database
	db, err := sqlx.Connect("pgx", cfg.Postgres.ConnString)
	if err != nil {
		l.Error(err)
		return
	}

	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			l.Error(err)
			return
		}
	}(db)

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Postgres.ConnMaxIdleTime * time.Second)

	err = db.Ping()
	if err != nil {
		l.Error(err)
		return
	}

	// auto-apply migrations if configured
	if cfg.Postgres.AutoMigrate {
		migrationDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			l.Error(err)
			return
		}

		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file://%s", cfg.Postgres.MigrationsPath),
			"user",
			migrationDriver,
		)
		if err != nil {
			l.Error(err)
			return
		}

		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			l.Error(err)
			return
		}
	}

	l.Debug("Connected to PostgreSQL")

	// creating redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// redis client connection check
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil || pong != "PONG" {
		l.Errorf("unable to ping redis: %v", err)
		return
	}

	l.Infof("connected to redis successfully")

	// creating fiber app
	f := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          v1.HandleError(),
	})
	f.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	// creating repositories
	userRepo := gateway.NewUserRepository(db, redisClient)

	// creating services
	userService := service.NewUserService(userRepo)

	// creating controllers
	redisController := v1.NewRedisController(userService)
	postgresController := v1.NewPostgresController(userService)
	signController := v1.NewSignController(userService)

	// defining groups
	redisGroup := f.Group("redis")
	postgresGroup := f.Group("postgres")
	signGroup := f.Group("sign")

	// registering http routes
	redisController.RegisterRoutes(redisGroup)
	postgresController.RegisterRoutes(postgresGroup)
	signController.RegisterRoutes(signGroup)

	go func() {
		err = f.Listen(net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port))
		if err != nil {
			l.Fatal(err.Error())
		}
	}()

	l.Debug("Started HTTP server")

	l.Debug("Application has started")

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

	l.Info("Application has been shut down")
}
