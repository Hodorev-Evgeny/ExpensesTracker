package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	core_pgx_pool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
	core_middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	core_transport_server "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/category/repository"
	feature_service_categor "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/category/service"
	feature_category_transport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/category/transport/http"
	feature_repository_limit "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/limit/repository"
	feature_service_limit "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/limit/service"
	feature_transport_limit "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/limit/transport"
	feature_repository_static "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/static/repository"
	feature_service_static "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/static/service"
	feature_transport_static "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/static/transport"
	feature_repository_transaction "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/repository"
	feature_transaction_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/service"
	feature_transactio_transport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/transport"
	features_users_repository "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/repository/postgres"
	feature_user_service "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/service"
	features_users_transport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/transport/http"
	feature_repository_file_system "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/web/repository/file_system"
	feature_service_web "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/web/service"
	feature_transport_web "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/web/transport"
	"go.uber.org/zap"

	_ "github.com/Hodorev-Evgeny/ExpensesTracker/docs"
)

// @title ExpensesTracker
// @version 0.7
// @description This server help you tracker money
// @host 127.0.0.1:8080
// @BasePath /api/v1
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()
	fmt.Println("starting server")

	serverConfig := core_transport_server.MustNewConfigServer()
	time.Local = serverConfig.TimeZone

	config := core_logger.MustNewConfig()
	logger, err := core_logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Error initializing logger: %v", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("starting expenses app")

	logger.Debug("starting initialization pool connection")
	pgconfig := core_pgx_pool.MustPostgresConfig()
	pool := core_pgx_pool.CreatePoolMust(ctx, pgconfig)
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("error pinging pool", zap.Error(err))
		os.Exit(1)
	}

	logger.Debug("starting initialization user service")
	userRepo := features_users_repository.NewUserRepository(pool)
	userServ := feature_user_service.NewUserService(userRepo)

	logger.Debug("starting initialization user transport")
	userTransporthttp := features_users_transport.NewUserHTTPHandler(userServ)
	userRouters := userTransporthttp.Routers()

	categoryRepository := feature_repositor_category.NewCategoryRepository(pool)
	categoryService := feature_service_categor.NewCategoryService(categoryRepository)
	categoryTransport := feature_category_transport.NewCategoryHTTPHandler(categoryService)
	categoryRouters := categoryTransport.Routes()

	transactionRepository := feature_repository_transaction.NewTransactionRepository(pool)
	transactionService := feature_transaction_service.NewTransactionService(transactionRepository)
	transactionTransport := feature_transactio_transport.NewTransactionHTTPHandler(transactionService)
	transactionRouters := transactionTransport.Router()

	limitRepository := feature_repository_limit.NewLimitRepository(pool)
	limitService := feature_service_limit.NewLimitService(limitRepository)
	limitTransport := feature_transport_limit.NewLimitHTTPHandler(limitService)
	limitRouters := limitTransport.Router()

	staticRepository := feature_repository_static.NewStaticRepository(pool)
	staticService := feature_service_static.NewStaticService(staticRepository)
	staticTransport := feature_transport_static.NewStaticHTTPHandler(staticService)
	staticRouters := staticTransport.Router()

	fileRepository := feature_repository_file_system.NewWebRepository()
	webService := feature_service_web.NewWebService(fileRepository)
	webTransport := feature_transport_web.NewWebTransport(webService)
	webRouters := webTransport.Router()

	apiVersionRouter := core_transport_server.NewAPIVersionRouter(core_transport_server.ApiVersion1)
	apiVersionRouter.RegisterAPIRoutes(userRouters...)
	apiVersionRouter.RegisterAPIRoutes(categoryRouters...)
	apiVersionRouter.RegisterAPIRoutes(transactionRouters...)
	apiVersionRouter.RegisterAPIRoutes(limitRouters...)
	apiVersionRouter.RegisterAPIRoutes(staticRouters...)

	httpServer := core_transport_server.NewServer(
		serverConfig,
		logger,
		core_middleware.CORS(serverConfig.AllowedOrigins),
		core_middleware.RequestId(),
		core_middleware.Logger(logger),
		core_middleware.Trace(),
		core_middleware.PanicRecovery(),
	)

	httpServer.ResisterApiVersionRouter(apiVersionRouter)
	httpServer.AddFrond()
	httpServer.RegisterRoutes(webRouters...)
	httpServer.RegisterSwagger()

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server failed to start", zap.Error(err))
	}
}
