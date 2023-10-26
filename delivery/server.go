package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mnc-test/config"
	"mnc-test/delivery/controller"
	"mnc-test/delivery/middleware"
	"mnc-test/manager"
)

type Server struct {
	um          manager.UsecaseManager
	gin         *gin.Engine
	log         *zap.Logger
	host        string
	redisClient *redis.Client
}

// middleware goes here
func (s *Server) InitMiddleware() {
	s.gin.Use(middleware.ZapLogger(s.log))
}

// controller
func (s *Server) InitController() {
	controller.NewUserCredentialController(s.um.UserCredUsecase(), s.gin, s.redisClient, s.log).Route()
	controller.NewMerchantController(s.um.MerchantUsecase(), s.gin, s.redisClient, s.log).Route()
	controller.NewTransactionController(s.um.TransactionUsecase(), s.gin, s.redisClient, s.log).Route()
	controller.NewTransferController(s.um.TransferUsecase(), s.gin, s.redisClient, s.log).Route()
}

// run server
func (s *Server) Run() {
	s.InitMiddleware()
	s.InitController()
	err := s.gin.Run(s.host)
	if err != nil {
		s.log.Error("Failed to run server", zap.String("error", err.Error()))
	}
}

func NewServer() (*Server, error) {
	// zap logger goes here baby
	zapconf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"activity.log"}, //activity from merchants,customer,etc goes here. we don't want to mix the log output to log.log
		ErrorOutputPaths: []string{"activity.log"}, //sameeee
	}
	logger, err := zapconf.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("Failed on config server", zap.String("error", err.Error()))
		return nil, err
	}

	//constructor from infra
	im, err := manager.NewInfraManager(cfg)
	if err != nil {
		logger.Error("Failed on construct infra", zap.String("error", err.Error()))
		return nil, err
	}

	//get the Redis client from the infraManager
	redisClient := im.GetRedisClient()

	//constructor from repomanager
	rm := manager.NewRepoManager(im)
	//contructor from usecase manager
	um := manager.NewUsecaseManager(rm, redisClient)

	//set host for gin server
	host := fmt.Sprintf("%s:%s", cfg.ApiConfig.Host, cfg.ApiConfig.Port)
	//return gin instance
	g := gin.Default()

	if err != nil {
		return nil, fmt.Errorf("Failed to construct infra: %v", err)
	}

	return &Server{
		um:          um,
		gin:         g,
		log:         logger,
		host:        host,
		redisClient: redisClient,
	}, nil
}
