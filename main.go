package main

import (
	"context"
	"time"
	"user-management-service/cmd/http"
	"user-management-service/cmd/middleware"
	"user-management-service/common"
	"user-management-service/common/logger"
	"user-management-service/config"
	db "user-management-service/config/database"
	"user-management-service/delivery"
	user3 "user-management-service/delivery/user"
	"user-management-service/repository"
	"user-management-service/repository/user"
	"user-management-service/service"
	user2 "user-management-service/service/user"
)

func main() {
	ctx := context.Background()
	// Start Init //
	loc, err := time.LoadLocation(config.LoadTimeZoneFromEnv())
	if err != nil {
		panic(err)
	}
	time.Local = loc
	// Configuration
	config.Init()

	// Logger
	logger.Init(logger.Config{
		AppName: config.Cold.AppName,
		Debug:   config.Hot.AppDebug,
	})

	database, err := db.NewDB(&db.Config{
		Driver:                config.Cold.DBMysqlDriver,
		Host:                  config.Cold.DBMysqlHost,
		Port:                  config.Cold.DBMysqlPort,
		DBName:                config.Cold.DBMysqlDBName,
		User:                  config.Cold.DBMysqlUser,
		Password:              config.Cold.DBMysqlPassword,
		SSLMode:               config.Cold.DBMysqlSSLMode,
		MaxOpenConnections:    config.Cold.DBMysqlMaxOpenConnections,
		MaxLifeTimeConnection: config.Cold.DBMysqlMaxLifeTimeConnection,
		MaxIdleConnections:    config.Cold.DBMysqlMaxIdleConnections,
		MaxIdleTimeConnection: config.Cold.DBMysqlMaxIdleTimeConnection,
	})
	if err != nil {
		panic(err)
	}

	// Registry
	commonValidator := common.NewValidator()
	commonRegistry := common.NewRegistry(
		common.WithValidator(commonValidator),
	)
	// End Init //

	// Start Repositories //
	userRepository := user.NewUserRepository(commonRegistry, database)
	repositories := repository.NewRegistry(database, userRepository)
	// End Repositories //

	// Start Services //
	userService := user2.NewUserService(commonRegistry, repositories)
	services := service.NewRegistry(
		userService,
	)
	// End Services //

	// Start Deliveries //
	userDelivery := user3.NewUserDelivery(commonRegistry, services)
	deliveries := delivery.NewRegistry(
		userDelivery,
	)
	// End Deliveries //

	// Start Middleware
	middlewares := middleware.NewMiddleware(commonRegistry, services)
	// End Middleware
	// Start HTTP Server //
	httpServer := http.NewServer(
		commonRegistry,
		deliveries,
		middlewares,
	)

	httpServer.Serve(ctx)
	// End HTTP Server //
}
