package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hodorev-Evgeny/ExpensesTracker"
	_ "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	corelogger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	pgpool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/postgresql/pool/pgx"
	redispool "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/repository/redis/pool/goredis"
	middleware "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/middleware"
	httpserver "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/server"
	categoryrepo "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/category/repository"
	categoryservice "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/category/service"
	categorytransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/category/transport/http"
	healthrepo "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/health/repository"
	healthtransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/health/transport"
	limitrepo "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/limit/repository"
	limitservice "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/limit/service"
	limittransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/limit/transport"
	statsrepo "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/static/repository"
	statsservice "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/static/service"
	statstransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/static/transport"
	transactionrepo "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/repository"
	transactionservice "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/service"
	transactiontransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/transaction/transport"
	userrepo "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/repository/postgres"
	userredis "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/repository/redis"
	userservice "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/service"
	usertransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/users/transport/http"
	webfiles "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/web/repository/file_system"
	webservice "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/web/service"
	webtransport "github.com/Hodorev-Evgeny/ExpensesTracker/internal/features/web/transport"
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

	serverConfig := httpserver.MustNewConfigServer()
	time.Local = serverConfig.TimeZone

	config := corelogger.MustNewConfig()
	logger, err := corelogger.NewLogger(config)
	if err != nil {
		fmt.Printf("Error initializing logger: %v", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("starting expenses app")

	logger.Debug("starting initialization pool connection")
	pgconfig := ExpensesTracker.MustPostgresConfig()
	pool := pgpool.CreatePoolMust(ctx, pgconfig)
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("error pinging pool", zap.Error(err))
		os.Exit(1)
	}

	logger.Debug("starting initialization redis connection")
	redisConfig := redispool.MustGetRedisConfig()
	redisClient := redispool.CreateRedisClientMust(redisConfig)
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("error closing redis client", zap.Error(err))
		}
	}()

	ctxCheckRedis, cancelRedis := context.WithTimeout(ctx, redisConfig.Timeout)
	defer cancelRedis()

	_, errRedisPing := redisClient.Client.Ping(ctxCheckRedis).Result()
	if errRedisPing != nil {
		logger.Error("error pinging redis", zap.Error(errRedisPing))
		os.Exit(1)
	}

	userRedis := userredis.NewRepositoryRedis(redisClient)

	logger.Debug("starting initialization user service")

	userRepo := userrepo.NewUserRepository(pool)
	userServ := userservice.NewUserService(userRepo, userRedis)

	logger.Debug("starting initialization user transport")
	userTransporthttp := usertransport.NewUserHTTPHandler(userServ, serverConfig.Secure, redisConfig.TTL)
	userRouters := userTransporthttp.Routers()

	categoryRepository := categoryrepo.NewCategoryRepository(pool)
	categoryService := categoryservice.NewCategoryService(categoryRepository)
	categoryTransport := categorytransport.NewCategoryHTTPHandler(categoryService)
	categoryRouters := categoryTransport.Routes()

	transactionRepository := transactionrepo.NewTransactionRepository(pool)
	transactionService := transactionservice.NewTransactionService(transactionRepository)
	transactionTransport := transactiontransport.NewTransactionHTTPHandler(transactionService)
	transactionRouters := transactionTransport.Router()

	limitRepository := limitrepo.NewLimitRepository(pool)
	limitService := limitservice.NewLimitService(limitRepository)
	limitTransport := limittransport.NewLimitHTTPHandler(limitService)
	limitRouters := limitTransport.Router()

	staticRepository := statsrepo.NewStaticRepository(pool)
	staticService := statsservice.NewStaticService(staticRepository)
	staticTransport := statstransport.NewStaticHTTPHandler(staticService)
	staticRouters := staticTransport.Router()

	fileRepository := webfiles.NewWebRepository()
	webService := webservice.NewWebService(fileRepository)
	webTransport := webtransport.NewWebTransport(webService)
	webRouters := webTransport.Router()

	healthRepository := healthrepo.NewHealthRepository(redisClient, pool)
	healthTransport := healthtransport.NewHealthTransport(healthRepository)
	healthRouters := healthTransport.Router()

	apiVersionRouter := httpserver.NewAPIVersionRouter(httpserver.ApiVersion1)
	apiVersionRouter.RegisterAPIRoutes(userRouters...)
	apiVersionRouter.RegisterAPIRoutes(categoryRouters...)
	apiVersionRouter.RegisterAPIRoutes(transactionRouters...)
	apiVersionRouter.RegisterAPIRoutes(limitRouters...)
	apiVersionRouter.RegisterAPIRoutes(staticRouters...)
	apiVersionRouter.RegisterAPIRoutes(healthRouters...)

	httpServer := httpserver.NewServer(
		serverConfig,
		logger,
		middleware.CORS(serverConfig.AllowedOrigins),
		middleware.RequestId(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.PanicRecovery(),
		middleware.Authenticator(redisClient),
	)

	httpServer.ResisterApiVersionRouter(apiVersionRouter)
	httpServer.AddFrond()
	httpServer.RegisterRoutes(webRouters...)
	httpServer.RegisterSwagger()

	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server failed to start", zap.Error(err))
		os.Exit(1)
	}
}
